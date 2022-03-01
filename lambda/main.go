package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jeremytorres/rawparser"
	"github.com/nf/cr2"
	"github.com/nfnt/resize"
)

func resizeImageThumbnail(ctx context.Context, img image.Image, width uint, height uint) (res io.Reader, err error) {
	m := resize.Resize(width, height, img, resize.NearestNeighbor)
	res, w := io.Pipe()
	go func() {
		defer w.Close()
		err = jpeg.Encode(w, m, nil)
		fmt.Println("encoded", err)
		if err != nil {
			panic(err)
		}
	}()
	return
}

func decodeNefImage(data []byte) (img image.Image, err error) {
	tmpDir := os.TempDir()
	sourceFile := path.Join(tmpDir, "file.nef")
	os.WriteFile(sourceFile, data, 0777)
	parser, _ := rawparser.NewNefParser(true)
	info := &rawparser.RawFileInfo{
		File:    sourceFile,
		Quality: 100,
		DestDir: tmpDir + "/",
	}
	file, err := parser.ProcessFile(info)
	if err != nil {
		return
	}
	resultData, err := os.ReadFile(file.JpegPath)
	if err != nil {
		return
	}

	img, _, err = image.Decode(bytes.NewReader(resultData))
	return
}
func getIconImage(filename string) (img image.Image, err error) {
	data, err := os.ReadFile("./icons/" + filename)
	if err != nil {
		return
	}
	img, _, err = image.Decode(bytes.NewReader(data))
	return
}

func decodeImage(contentType string, data []byte) (img image.Image, err error) {
	fmt.Println("object content type:", contentType)
	img, err = decodeNefImage(data)
	if err == nil {
		return
	}
	img, _, err = image.Decode(bytes.NewReader(data))
	if err == nil {
		return
	}
	img, err = cr2.Decode(bytes.NewReader(data))
	if err == nil {
		return
	}
	img, err = getIconImage("file.png")
	return
}

func generateThumbnailFromS3File(ctx context.Context, svc *s3.S3, svcUpload *s3manager.Uploader, bucket string, key string, thumbnail string, width uint, height uint) (err error) {
	resp, err := svc.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	fmt.Println("received object", err)
	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	fmt.Println("read all object data", len(data), err)
	if err != nil {
		return
	}

	contentType := *resp.ContentType
	fmt.Println("object content type:", contentType)
	img, err := decodeImage(contentType, data)
	if err != nil {
		return
	}

	r, err := resizeImageThumbnail(ctx, img, width, height)
	if err != nil {
		return
	}

	sourceBucket := os.Getenv("BUCKET")
	_, err = svcUpload.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket:      aws.String(sourceBucket),
		Key:         aws.String(thumbnail),
		Body:        r,
		ContentType: aws.String("image/jpeg"),
	})
	fmt.Println("uploaded", sourceBucket, "/", thumbnail, err)

	return
}

func handleRequest(svc *s3.S3, svcUpload *s3manager.Uploader) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// bucket := os.Getenv("BUCKET")
	redirectURL := os.Getenv("REDIRECT_URL")
	// privKeyData, _ := base64.StdEncoding.DecodeString("LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBbndVWEEzMEs4WGhleFVYelhKRWh6Njh2YkhpdE1IN1NLQnE5azg3OEUzRkRLRnA2CkkrcXl2bzBySjBrNTJoWDFMN2IxNHVVZjdJeTJHMlVydjdzZkFmR1hVVDMyWnE4SUlBeEE2VXJKUXJ3Q2NsV1AKdllMVHVZdGVoOHBVeFdtb0tQZWt4M3ZyaTIrcHNyYm5XNkFVRkdVZy9QUjZQL0FPM3BaaTc2OW5DZkJzS3l2RwpIRWtqck1qK0FidUJMUjVTbzI5NHVya1JHc00rQlZtTnBMZmVUWkZrdXZPY2hMM2tOTnluaWtxNXZuVEd5S3MrCjdpcUJOU0xMdDErc0lHWENUV3lLREMvbTRuODBCTyt6dzJubXJoeUZibnZnTXhtU1A0MG5jcnRaWjZ0TDVYcnEKeldPRldXWnVBdlpBNzBtYlF2VHp3RFBib210a0I1a2JwRzFnS1FJREFRQUJBb0lCQURtTjlnNWxLNWpLSUVMMgpWbWFpZ01KU2ZhK2MzZEkzbElpL1pPNGlkRW4xTElGbWZkNlNDNis4R0dwWHJvQ29nZDJMTXVPMFdpR2JyQzJ1CktHOTBvbnhwbklMakVsT0g5a0RqTG50QmZpYVJOYkY4RmxKWkQ1aVFRTjZVOUttZTNDWTN1NFFYN2VEQndZSFYKZ1ZkQVVRUXd0Z2ZYMUlkZnM0SU1uREowSWE4T3EyUEdqVVp6c0hwUnBOWWJnZWpUMHBKenUyNWE4NXVMVDJtNQpNL0I3cDg3MjRlbllmRkEvY2RuY3YrY2FhT0RVcjNDRGdXQTBZRWZnVHNJVlpKQUtMM3ZWK2FpOWlpbDJSbmg2CmoyenMzdi8vc21JTEFMcnY4RDNoSExJSUtnV2VoeXNJZkd1S2xrM0RnbkE1SnFHeUs1STQ1bVZaNmdYejc5K0QKNnJibHdLRUNnWUVBMHlGS3pWd2lBMnA2TmkyeWlad2ROcEFKNDRyZEZkRngvdStrV0xvVlVHR3NMZGZlSjhFVwpVUS9zUC9yQVZsQ1l5M2ZsMU9ZWm1yL3cvL2FTNVBEV0pwZ1ZPc3RIT04xNmoyN0FkVkVjNm5SYTVadHBGVEQrCm42RGdLdHF1SUs2cGdGeEVVZkVYUlhhK0pPM1JYZzZ1eUNKN21kWURWd3pVTXh3V0FWY3ZQTU1DZ1lFQXdOQ3oKU1VDZUIxUDhFMHJzSVF5cktmN0htaXdUNkwxSTF1ZFdGYVNkK2FWMG5lSXh5N09OL3N3RitId3M0YlhreVJoeAovUjJzNlp0RHJ1TzVHMlZIUW9veXVtcUlvenpnTlRoK2ExbGpaRGFwZ25NdFd5SUxTU2dkakE3ajZERFhoZGVkClF4U3dZb0VTcDZ2eFJGTkJxaEZRNytoSlpBa016MzJQUUxGa2tLTUNnWUVBenc0SFhmN05IS3gvemtKaXBiSUUKdEUzdVpNajZxVHNPb1FaRUZ2L25oejFDcm5MVjNBYnc2KzdCc3IwbmczN25XaVByc2k4M1RSeVFMWGFUK1JKMwo4c0dUa0dWckk4bVJPTGxVNWJqMnNyZ0pyTVFNK2t0aWF3cEt6YnhJcUtTaWR0QUY4SmFRUy85MzJwK2doSzNCCm0yUzE4dGgvemc4MnpDanZLOEZsQWlFQ2dZRUFzbVBrbkU5V1pmMUQ2UzJXVGZXRW52UUVCQlhuelpyaXUwR0oKR2JrV0Y1VUcwZFFtc2dwTHc0Tkx1dHhZUWZPaDJwUHRVbnVVTVFYZmx2MUZrNTBlVXVlOWkwOXBYMjNCR2p4TQphbEZuYlo1Tk1rNFJscEZtMDZaenY5TSs5T0hMWlI5WmRtaTcwRWNPMVdaMWIvdC9jek5XS01CR3RuRFJFMTlkCm5FTURnZlVDZ1lBdEgxZnZRL2VJWnhSb0hPcjI4bWV3d2dOT3R0bjdHSmVwRHBVK2cyRUdsKzBDN2xaUDUyWlMKRC9hVEsyZ3JRR0lsM0lqaHd4WXo3bFIzUGhGMWplcFRjT0ZscHVndXFqSklOYVN5VmRJeVhwaXVsQWRGR0FjbQpiWnhvRlNDM2tjcVpnWCtQWm9pV242d3hEVFVjeFoxbUJJSXFYcm1tVGV6UERiYU1rK0xRUlE9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=")
	// fmt.Println("??", string(privKeyData))
	// privKey, _ := sign.LoadPEMPrivKey(bytes.NewReader(privKeyData))
	// signer := sign.NewURLSigner("K2HLPDDF1Z0Z8R", privKey)
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {
		thumbnail := request.Path
		parts := strings.Split(request.Path, "-")
		if len(parts) > 1 {
			key := strings.Join(parts[:len(parts)-1], "-")
			var width uint64
			width, err = strconv.ParseUint(parts[len(parts)-1], 0, 32)
			if err != nil {
				return
			}
			if width > 2000 {
				err = fmt.Errorf("width %d is greater then maximum allowed (2000)", width)
				return
			}
			height := uint(0)
			keyParts := strings.Split(strings.TrimLeft(key, "/"), "/")
			bucket := keyParts[0]
			fmt.Println("path", key, "=>", bucket, "+", strings.Join(keyParts[1:], "/"), "key parts:", keyParts)
			key = strings.Join(keyParts[1:], "/")
			err = generateThumbnailFromS3File(ctx, svc, svcUpload, bucket, key, thumbnail, uint(width), height)
			if err != nil {
				return
			}
		}
		query := url.Values{}
		for key, value := range request.QueryStringParameters {
			query.Add(key, value)
		}
		rawURL, err := url.Parse(redirectURL)
		if err != nil {
			return
		}
		rawURL.Path = path.Join(rawURL.Path, thumbnail)
		// signedURL, err := signer.Sign(rawURL.String(), time.Now().Add(1*time.Hour))
		// if err != nil {
		// 	log.Fatalf("Failed to sign url, err: %s\n", err.Error())
		// }
		signedURL := rawURL.String()

		res = events.APIGatewayProxyResponse{
			Headers: aws.StringValueMap(map[string]*string{
				"location": aws.String(signedURL),
			}),
			StatusCode: 301,
		}
		return
	}
}

func main() {
	sess, sErr := session.NewSession(&aws.Config{
		// Region: aws.String("us-west-2")},
	})
	sess = session.Must(sess, sErr)
	svc := s3.New(sess)
	svcUpload := s3manager.NewUploader(sess)
	lambda.Start(handleRequest(svc, svcUpload))
}
