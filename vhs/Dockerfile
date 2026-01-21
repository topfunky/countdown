FROM ghcr.io/charmbracelet/vhs

# Install Dependencies
RUN apt-get update --allow-releaseinfo-change || true
RUN apt-get -y install build-essential procps curl file git 

# Prepare Homebrew directory and install as vhs user
RUN mkdir -p /home/linuxbrew/.linuxbrew /home/vhs/.cache && chown -R 1976:1976 /home/linuxbrew /home/vhs
RUN NONINTERACTIVE=1 HOME=/vhs su vhs -s /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

RUN echo >> /vhs/.bashrc
RUN echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv bash)"' >> /vhs/.bashrc

# Ensure vhs user can use homebrew
RUN chown -R 1976:1976 /home/linuxbrew/.linuxbrew

ENV PATH="/home/linuxbrew/.linuxbrew/bin:${PATH}"

# Install packages
RUN su vhs -s /bin/bash -c "/home/linuxbrew/.linuxbrew/bin/brew install go" 

