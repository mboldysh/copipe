# Copipe

Copipe

![Demonstration](copipe-demo-2.gif)

## Installation

Cuurently, installation is only possible using binaries. Use following steps for installation:

1. Open [Releases Page](https://github.com/mboldysh/copipe/releases)
2. Select release version that you want to install
3. Open your terminal and execute commands below

```bash
# Specify here version that you want to install
export COPIPE_VERSION=0.1.0

# Donwload and install
curl -fSL "https://github.com/mboldysh/copipe/releases/download/${COPIPE_VERSION}/copipe_${COPIPE_VERSION}_darwin_amd64.tar.gz" -o copipe.tar.gz \
  && tar -xf copipe.tar.gz \
  && mv copipe /usr/local/bin/copipe \
  && chmod a+x "/usr/local/bin/copipe" \
  && rm -f copipe.tar.gz
```

## Usage