FROM alpine:latest
WORKDIR /app
RUN addgroup -S app && \
    adduser -S app -G app
COPY build/app /app
RUN chown -R app:app /app && ls -lah /app/* && chmod 700 /app/app
USER app

ENTRYPOINT [ "/app/app", "consumer" ]