class Ccm < Formula
  desc "Claude Code Manager - Manage multiple AI model providers for Claude Code"
  homepage "https://github.com/taliove/go-claude-model"
  version "0.1.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/taliove/go-claude-model/releases/download/v0.1.0/ccm-darwin-arm64"
      sha256 "YOUR_SHA256_ARM64"
      bottle :unneeded
    elsif Hardware::CPU.intel?
      url "https://github.com/taliove/go-claude-model/releases/download/v0.1.0/ccm-darwin-amd64"
      sha256 "YOUR_SHA256_AMD64"
      bottle :unneeded
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.arm?
      url "https://github.com/taliove/go-claude-model/releases/download/v0.1.0/ccm-linux-arm64"
      sha256 "YOUR_SHA256_LINUX_ARM64"
      bottle :unneeded
    elsif Hardware::CPU.intel?
      url "https://github.com/taliove/go-claude-model/releases/download/v0.1.0/ccm-linux-amd64"
      sha256 "YOUR_SHA256_LINUX_AMD64"
      bottle :unneeded
    end
  end

  def install
    bin.install "ccm"
  end

  def post_install
    ohai "CCM installed successfully!"
    puts "Run 'ccm init' to get started"
    puts "Or 'ccm --help' for usage information"
  end

  test do
    system "#{bin}/ccm", "--version"
  end
end
