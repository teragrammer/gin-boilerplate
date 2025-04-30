FROM ubuntu:24.04

# Set environment variables to ensure non-interactive installation (avoid prompts)
ENV DEBIAN_FRONTEND=noninteractive

# Update the package index and install dependencies for Go, Python, AWS CLI, and some security tools
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y \
    build-essential \
    libssl-dev \
    zlib1g-dev \
    libbz2-dev \
    libreadline-dev \
    libsqlite3-dev \
    wget \
    llvm \
    libncurses5-dev \
    libgdbm-dev \
    libnss3-dev \
    libffi-dev \
    liblzma-dev \
    tk-dev \
    libdb-dev \
    libexpat1-dev \
    libpcap-dev \
    make \
    gcc \
    git \
    ca-certificates \
    curl \
    gnupg \
    lsb-release \
    software-properties-common \
    build-essential \
    unzip \
    python3-venv \
    jq

# Set the version of Python to install
ENV PYTHON_VERSION=3.13.3

# Download and install Python from source
RUN wget https://www.python.org/ftp/python/${PYTHON_VERSION}/Python-${PYTHON_VERSION}.tgz && \
    tar xzf Python-${PYTHON_VERSION}.tgz && \
    cd Python-${PYTHON_VERSION} && \
    ./configure --enable-optimizations && \
    make -j$(nproc) && \
    make altinstall

# Create a virtual environment to avoid 'externally-managed-environment' issue
RUN python3 -m venv /venv

# Activate the virtual environment and upgrade pip inside it
RUN /venv/bin/pip install --no-cache-dir --upgrade pip

# Set the environment variable to use the virtual environment by default
ENV PATH="/venv/bin:$PATH"

# Install the latest AWS CLI
RUN pip install awscli --upgrade
RUN pip install awscli-local --upgrade

# Install Go
RUN wget https://go.dev/dl/go1.24.2.linux-amd64.tar.gz
RUN tar -C /usr/local -xvf go1.24.2.linux-amd64.tar.gz
RUN echo "export PATH=$PATH:/usr/local/go/bin:~/.local/bin:$PATH" >> ~/.bashrc

# NODEJS
RUN curl -fsSL https://deb.nodesource.com/setup_22.x | bash - && apt-get update && apt-get install -y nodejs && rm -rf /var/lib/apt/lists/*
# Verify that Node.js and npm were installed correctly
RUN node -v
RUN npm -v

# Update npm packages
RUN npm install -g npm@latest
RUN npm install -g nodemon
RUN npm install -g npm-check-updates

# Clean up unnecessary files and clear the apt cache to reduce image size
RUN apt-get autoremove -y && apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Set the working directory
WORKDIR /app

# Expose port
EXPOSE $APP_PORT
