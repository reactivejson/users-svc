FROM scratch

ENV APP=/usr/local/bin/user-svc \
    USER_UID=1101 \
    USER_NAME=user-svc

# install engine binary
COPY build/_output/bin/user-svc ${APP}
USER ${USER_UID}:${USER_UID}
CMD [ "./${APP}" ]
