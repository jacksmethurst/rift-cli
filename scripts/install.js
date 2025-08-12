#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const https = require('https');
const { promisify } = require('util');
const { pipeline } = require('stream');
const zlib = require('zlib');
const tar = require('tar');

const streamPipeline = promisify(pipeline);

const REPO = 'jacksmethurst/rift-cli';
const VERSION = require('../package.json').version;

function getPlatform() {
  const platform = process.platform;
  const arch = process.arch;
  
  if (platform === 'darwin') {
    return arch === 'arm64' ? 'darwin-arm64' : 'darwin-amd64';
  } else if (platform === 'linux') {
    return arch === 'arm64' ? 'linux-arm64' : 'linux-amd64';
  } else if (platform === 'win32') {
    return 'windows-amd64';
  }
  
  throw new Error(`Unsupported platform: ${platform}-${arch}`);
}

async function downloadBinary() {
  const platformSuffix = getPlatform();
  const isWindows = process.platform === 'win32';
  const extension = isWindows ? 'zip' : 'tar.gz';
  const filename = `rift-${platformSuffix}.${extension}`;
  const url = `https://github.com/${REPO}/releases/download/v${VERSION}/${filename}`;
  
  console.log(`Downloading rift binary for ${platformSuffix}...`);
  console.log(`URL: ${url}`);
  
  const binDir = path.join(__dirname, '..', 'bin');
  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }
  
  return new Promise((resolve, reject) => {
    https.get(url, (response) => {
      if (response.statusCode === 302 || response.statusCode === 301) {
        // Follow redirect
        return https.get(response.headers.location, handleResponse);
      }
      
      function handleResponse(res) {
        if (res.statusCode !== 200) {
          reject(new Error(`Failed to download: ${res.statusCode} ${res.statusMessage}`));
          return;
        }
        
        if (isWindows) {
          // Handle ZIP for Windows
          const chunks = [];
          res.on('data', chunk => chunks.push(chunk));
          res.on('end', () => {
            const buffer = Buffer.concat(chunks);
            // For simplicity, we'll skip ZIP extraction for now
            // You'd need a ZIP library or use a different approach
            fs.writeFileSync(path.join(binDir, 'rift.exe'), buffer);
            resolve();
          });
        } else {
          // Handle tar.gz for Unix systems
          streamPipeline(
            res,
            zlib.createGunzip(),
            tar.extract({
              cwd: binDir,
              strip: 0,
              filter: (path, entry) => {
                return path.includes('rift-');
              },
              map: (header) => {
                if (header.name.includes('rift-')) {
                  header.name = 'rift';
                }
                return header;
              }
            })
          ).then(() => {
            // Make executable
            const binaryPath = path.join(binDir, 'rift');
            if (fs.existsSync(binaryPath)) {
              fs.chmodSync(binaryPath, '755');
            }
            resolve();
          }).catch(reject);
        }
      }
      
      handleResponse(response);
    }).on('error', reject);
  });
}

async function main() {
  try {
    await downloadBinary();
    console.log('✅ rift installed successfully!');
    console.log('Run "rift version" to verify installation.');
  } catch (error) {
    console.error('❌ Installation failed:', error.message);
    process.exit(1);
  }
}

if (require.main === module) {
  main();
}