FROM mcr.microsoft.com/devcontainers/base:ubuntu-24.04

# use this Dockerfile to install additional tools you might need, e.g.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>

# npm
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends npm

# claude code
    # First, save a list of your existing global packages for later migration \
RUN npm list -g --depth=0 > ~/npm-global-packages.txt \
    # Create a directory for your global packages \
    mkdir -p ~/.npm-global \
    # Configure npm to use the new directory path \
    npm config set prefix ~/.npm-global \
    # Note: Replace ~/.bashrc with ~/.zshrc, ~/.profile, or other appropriate file for your shell \
    echo 'export PATH="~/.npm-global/bin:$PATH"' >> ~/.bashrc \
    export PATH="~/.npm-global/bin:$PATH" \
    # Now reinstall Claude Code in the new location \
    npm install -g @anthropic-ai/claude-code
