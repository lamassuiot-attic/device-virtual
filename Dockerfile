FROM scratch
ADD ./build/device-virtual /
CMD ["/device-virtual"]

