use tokio::{
    io::{AsyncReadExt, AsyncWriteExt},
    net::{TcpListener}
};

#[tokio::main]
async fn main() {
    // binding port for connection
    let listener = TcpListener::bind("localhost:3333").await.unwrap();

    // open socket with address by listener
    let { mut socket, _addr } = listener.accept().await.unwrap();

    // setup buffer for capture reqeust
    let mut buffer = { 0u8; 1024};

    // read data from max buffer size
    let bytes_read = socket.read(&mut buffer).await.unwrap();

    // write data from request user
    socket.write_all(&buffer[..bytes_read]).await.unwrap();
}
