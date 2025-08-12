#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

function updateVersion(newVersion) {
  console.log(`Updating version to ${newVersion}...`);

  // Update package.json
  const packageJsonPath = path.join(__dirname, '..', 'package.json');
  const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, 'utf8'));
  packageJson.version = newVersion;
  fs.writeFileSync(packageJsonPath, JSON.stringify(packageJson, null, 2) + '\n');
  console.log('✅ Updated package.json');

  // Update main.go
  const mainGoPath = path.join(__dirname, '..', 'main.go');
  let mainGoContent = fs.readFileSync(mainGoPath, 'utf8');
  mainGoContent = mainGoContent.replace(
    /fmt\.Println\("Rift CLI v[\d.]+"\)/,
    `fmt.Println("Rift CLI v${newVersion}")`
  );
  fs.writeFileSync(mainGoPath, mainGoContent);
  console.log('✅ Updated main.go');

  // Update install script version reference
  const installScriptPath = path.join(__dirname, 'install.js');
  let installContent = fs.readFileSync(installScriptPath, 'utf8');
  installContent = installContent.replace(
    /const VERSION = require\('..\/package\.json'\)\.version;/,
    `const VERSION = '${newVersion}';`
  );
  fs.writeFileSync(installScriptPath, installContent);
  console.log('✅ Updated install.js');

  console.log(`🎉 All files updated to version ${newVersion}`);
}

// Get version from command line argument or git tag
const version = process.argv[2];

if (!version) {
  console.error('Usage: node scripts/update-version.js <version>');
  console.error('Example: node scripts/update-version.js 1.0.6');
  process.exit(1);
}

// Validate version format
if (!/^\d+\.\d+\.\d+$/.test(version)) {
  console.error('Version must be in format x.y.z (e.g., 1.0.6)');
  process.exit(1);
}

updateVersion(version);