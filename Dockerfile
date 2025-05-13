#################################
#	STAGE 1 - Backend	#
#################################
FROM golang:tip-alpine3.20 AS backend
WORKDIR /opt/backend

# Copy backend sources
COPY . .

# Install required Go Dependencies
RUN go install ./cmd

# Build the project
RUN go build -o strata ./cmd

#################################
#	STAGE 2 - Frontend	#
#################################
FROM node:current-slim AS frontend
WORKDIR /opt/frontend

# Copy UI sources
COPY ui .

# Install UI lib(s)
RUN npm install

# Build frontend
RUN npm run build


#################################
#	STAGE 3 - Production	#
#################################
FROM golang:tip-alpine3.20 AS production
WORKDIR /opt/strata

# Copy in the backend sources
COPY --from=backend /opt/backend .
# Also install the strata executable
COPY --from=backend /opt/backend/strata /bin/strata

# Now copy in the frontend sources on top
COPY --from=frontend /opt/frontend ui

# Run strata
# This has to have no '/' characters or it will/can error, so we have to install it above.  
#
# If ./strata exists, start that. Otherwise fall back to the version in /bin/strata, installed when the docker image was built.
CMD ["sh", "-c", "if [ -f strata ] ; then ./strata ; else strata ; fi"]
