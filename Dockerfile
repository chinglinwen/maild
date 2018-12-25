FROM alpine
WORKDIR /app
COPY maild /app/maild
#ENTRYPOINT /app/maild
CMD /app/maild -h
EXPOSE 3001
