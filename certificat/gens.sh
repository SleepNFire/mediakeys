    rm *.pem
    rm *.srl
    openssl req -x509 -newkey rsa:4096 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/CN=*" -addext "subjectAltName = DNS:localhost,DNS:impression-tracking"
    echo "CA's self signed certificat"
    openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/CN=impression-tracking" -addext "subjectAltName = DNS:localhost,DNS:impression-tracking"
    openssl x509 -req -in server-req.pem -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server.conf
    echo "Server's signed certificate"
    echo "Verify Server certificate"
    openssl verify -CAfile ca-cert.pem server-cert.pem
    openssl req -newkey rsa:4096 -nodes -keyout client-key.pem -out client-req.pem -subj "/CN=ad-serving" -addext "subjectAltName = DNS:localhost,DNS:ad-serving"
    openssl x509 -req -in client-req.pem -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-cert.pem -extfile client.conf
    echo "Verify certificate"
    openssl verify -CAfile ca-cert.pem client-cert.pem