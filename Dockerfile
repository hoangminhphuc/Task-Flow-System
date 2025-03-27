FROM alpine

# will be created if not exist
WORKDIR /app/
ADD ./app /app/

ENTRYPOINT [ "./app" ]




