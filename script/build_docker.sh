image=openim/open_im_enterprise:v1.0.2
chmod +x ./*.sh
./build_all_service.sh
cd ../
docker build -t $image . -f ./deploy.Dockerfile
docker push $image
echo "build ok"
