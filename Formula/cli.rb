class Cli < Formula
  desc "A Git alternative written in Go"
  homepage "https://github.com/jacksmethurst/rift-cli"
  version "1.0.0"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/jacksmethurst/rift-cli/releases/download/v1.0.0/rift-darwin-amd64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_FOR_INTEL"
    end
    if Hardware::CPU.arm?
      url "https://github.com/jacksmethurst/rift-cli/releases/download/v1.0.0/rift-darwin-arm64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_FOR_ARM"
    end
  end

  def install
    bin.install "rift-darwin-amd64" => "rift" if Hardware::CPU.intel?
    bin.install "rift-darwin-arm64" => "rift" if Hardware::CPU.arm?
  end

  test do
    system "#{bin}/rift", "version"
  end
end