# terraform-provider-mirror-with-s3

This is local web server code used for the sole purpose of directly using the Terraform provider stored in AWS S3.
> In fact, if your Terraform-provider(.zip) object on AWS S3 is publicly accessible, you don't need it.


1. 로딩할때 s3에서 현재 서빙가능한 목록확인
2. 해당 목록을 기반으로 라우터 셋팅하고 서버 운영
3. 요청받으면 다운로드
 - config 넣을 구멍은없나??
4. 없으면, 말구~
