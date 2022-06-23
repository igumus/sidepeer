FROM scratch

COPY ./sidepeer /sidepeer

ENTRYPOINT [ "/sidepeer" ] 
