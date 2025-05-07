# Based on debian:latest
FROM debian
WORKDIR /opt/strata

# Expose the port 8080
# Expose to other containers, not sure if this is needed
#EXPOSE 8080

# Update and fetch packages
RUN apt-get -y update && apt-get -y upgrade
RUN apt-get install -y --no-install-recommends ca-certificates wget

# Download GO
RUN wget -O /tmp/go.tar.gz https://go.dev/dl/go1.24.1.linux-amd64.tar.gz

# Install GO 
RUN tar -C /usr/local -xzvf /tmp/go.tar.gz && rm /tmp/go.tar.gz
RUN ln -s /usr/local/go/bin/* /bin/.


# Lower the privilege level in case of attacks
RUN useradd -m strata
RUN chown -R strata:strata /opt/strata
USER strata

# Copy project contents into container
ADD --chown=strata . /opt/strata 

# Install required Go Dependencies
RUN go install ./cmd

# Build the project
RUN go build -o strata ./cmd

# Install
USER root
RUN cp strata /bin/strata
USER strata


# Run strata
# This has to have no '/' characters or it will/can error, so we have to install it above. 
# 
# If ./strata exists, start that. Otherwise fall back to the version in /bin/strata, installed when the docker image was built. 
CMD ["bash", "-c", "if [ -f strata ] ; then ./strata ; else strata ; fi"]
