# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.224.2/containers/go/.devcontainer/base.Dockerfile

# [Choice] Go version (use -bullseye variants on local arm64/Apple Silicon): 1, 1.16, 1.17, 1-bullseye, 1.16-bullseye, 1.17-bullseye, 1-buster, 1.16-buster, 1.17-buster
ARG VARIANT="1.17"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT} AS base
RUN curl -sSLJ "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"\
    && unzip awscliv2.zip \
    && ./aws/install && rm -r ./aws && rm awscliv2.zip

FROM alpine:latest as tools
RUN apk add --no-cache curl unzip

FROM tools AS kubectl
RUN curl -sSLJO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" \
    && chmod +x ./kubectl

FROM tools AS fzf
ARG FZF_VERSION
RUN curl -sSLJ "https://github.com/junegunn/fzf/releases/download/${FZF_VERSION}/fzf-${FZF_VERSION}-linux_amd64.tar.gz" -o "fzf.tar.gz"\
    && tar zxf fzf.tar.gz \
    && chmod +x ./fzf    

FROM base as final
ARG KREW_PLUGINS
COPY --from=kubectl "kubectl" "/usr/local/bin/kubectl"
COPY --from=fzf "fzf" "/usr/local/bin/fzf"
USER vscode
ENV PATH="/home/vscode/.krew/bin:${PATH}"
RUN KREW=$(mktemp -d) \
    && mkdir -p ${KREW} \
    && curl -sSLJ "https://github.com/kubernetes-sigs/krew/releases/latest/download/krew-linux_amd64.tar.gz" -o "${KREW}/krew.tar.gz"\
    && tar zxf ${KREW}/krew.tar.gz -C ${KREW} \
    && ${KREW}/krew-linux_amd64 install krew \
    && rm -rf ${KREW} \
    && kubectl krew install $KREW_PLUGINS
RUN go install github.com/goreleaser/goreleaser@latest
RUN sed -i 's/plugins=(git)/plugins=(git kubectl aws)/g' ~/.zshrc
