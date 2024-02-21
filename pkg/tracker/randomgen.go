package tracker

import (
	"math/rand"
	"time"

	"github.com/galexrt/fivenet/query/fivenet/model"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/zap"
)

func (s *Manager) randomizeUserMarkers() {
	userIdentifiers := []string{
		"char1:034c716223e6b4e431db6de501b28f65d8560cb1",
		"char1:04cc59ac55ee0f4d0ace5c3b519f5df929b5c6cf",
		"char1:08e6be478c776b6727a8fb3061f1431eb715bd3f",
		"char1:0a763f4bfb3aa03902b5da5cae7682c0a8c9a186",
		"char1:0ff2f772f2527a0626cac48670cbc20ddbdc09fb",
		"char1:1419c5a584086a76dbddae6aed1051157596a02a",
		"char1:1b6d306d910c03740ccdbed6074abbf96751071a",
		"char1:1ca6092b49f7931ed96823a5434ec85dece281cb",
		"char1:1d2429c1d7b7aed2b65d95097668c8bf18f6156b",
		"char1:1dd99d435c3094b397cad60b41823200ba83b293",
		"char1:1eaaf652f2ad2b1cae75922156c6cac25c8bf107",
		"char1:2736cacf0659f2aab7ea13b8138b48dd3684789e",
		"char1:2db159960d3cbf1a9598670d9f75bbe11ecf28e8",
		"char1:2fdae7dc444d9a6d00a19ad0789beeab682bb3f1",
		"char1:33142683fa04e746f65c928558a4f76f1aabd8c1",
		"char1:34b11943424e30d4e9cfcfb345a15b430858092e",
		"char1:40361766e04bc8614d296cf1049c91850a562e99",
		"char1:41b5d94557aadf9cc3feeca438c9aeaf76e1a0a7",
		"char1:4924b2719b149a4d838639a72953fef6f4bdbf01",
		"char1:4fc579b6a9e9678fad9f2b325157d4e73f03f013",
		"char1:5229f3a7f858312e3965aa0d9aaa3246a8b40db3",
		"char1:5ff610ff8a24276cdeb05ebe8e826afbefdb88bd",
		"char1:60368551546ee3b8c6d6c3d6830865ce664fac79",
		"char1:656d634eadfd0190acc459e207b458e4fe0bd546",
		"char1:6bc98f8c893ebbc01201bb0211d36fdb97948e41",
		"char1:6ec5bfdcae436036a602b20e3dd0da0b1b312638",
		"char1:6fedb0c5d9b76c4d9ab07b10cb52334a235e496c",
		"char1:724a887a2ce91af3e84d40771b3275fecc35be23",
		"char1:7292c6a393668e6b986d9ccfe6ff24b6911ceae0",
		"char1:742b31901f550a8ec512c375b9cd0434fa6112b6",
		"char1:749ff80ef754c3f1123564ae001ed771acb706ec",
		"char1:76bc318047b00616cfdaab09dc4da7cc30a2d730",
		"char1:801f3125701ee3186c44959f31c54846518385f0",
		"char1:8369b091d49b5aa2fe061c654d4624f5ce5b89f1",
		"char1:87903dfb67ad3ff56bae5b599b9ac00eb9528fd6",
		"char1:8875ccd943dff254f6973894e78f126ef3629eae",
		"char1:8a63032736f03e05d51e3e5b8ec7b3957e068ed6",
		"char1:8d0bf57a5bf239a98c7736d76625d75fccac19ba",
		"char1:8de0ea2bcbf2123317f50a6d2e8ab05e5f23029b",
		"char1:8e09f9132439f8bd3fa19d864968d83f492c1bb2",
		"char1:8e6137315121727dc3aa8ab0894fbf0e356ba7c2",
		"char1:9b2ba99e352f173f6c4949bbbd2dded50c4db435",
		"char1:a79ab349fc91d526b6891968278e9d2f61374a59",
		"char1:a9dffec518ac9b046745fd9beee7e13ba98292db",
		"char1:b1164816ff009de9c07b7c843667d28567d88dcd",
		"char1:b854005062b48e4dc2537b1776032e8c8792fdf2",
		"char1:b92fa727ae254a0ca89b8fc14bd2599c6ef98b22",
		"char1:bb8d0a58a0683f0482d13c99e4587769c597c9ff",
		"char1:c49d4b59e823ef98f5a487767d2729771e3fea01",
		"char1:c4c6e6f876a2417b3060025773fb193b24666e8b",
		"char1:c834a0f52b6d5465643ad9c8c689d2b5abe9801f",
		"char1:c88f0186b9655b90c42c88cefeb771d78a31b329",
		"char1:cbfdda4258a3651a699c947d82203e3a89df080b",
		"char1:cd1b7c7aa9d65397f1ec9ee1bee23a4e9994b915",
		"char1:d097792ced6bd291393b9a11f471736a89fe5b2e",
		"char1:d2466c943fab37a2486da961ec155c508777209b",
		"char1:d47542a5fbf8218373c2dfa892fb22a18c88df85",
		"char1:d4b145bb77a128e66c0cd96618a9950236cf0a70",
		"char1:da16df147bd41b781ff721b872d4fc4a1b5c6875",
		"char1:e3dea6f2d4624b9fcd83a4a89749eac0fa5d2a21",
		"char1:ea500d1ef7117cfc4a6df5fd580348ebf9a51e9f",
		"char1:eef505578bf3ee13dff3446a36e88439db2e5f5d",
		"char1:f474be485f087903da51ad3ef54e054b5ca85bd2",
		"char1:f9107882499841617eaf31cacf3611e64c067c3b",
		"char1:fac50986e86134a6867687fc65b13d9b332c1d35",
		"char1:fd6798c8e8ad98c90b4c22f39933ca8673b1aeb0",
		"char2:03f5b7bfc3659f8b4c72a240a249eb354589e077",
		"char2:0b166805186047311175d4b94a9d0921fd90c5c2",
		"char2:1419c5a584086a76dbddae6aed1051157596a02a",
		"char2:14dfc4185115fc439692a86d2cc16f1bfbde27d2",
		"char2:3b9f0e9150ef266d53acfbeb186aa475c457b10f",
		"char2:406561925f009d89ff962bd6b1117f9f4c52d354",
		"char2:64ed36dd998173afe9cc7ca7f90fe7845e0c4338",
		"char2:709567769c50e824fd2791c3147457c8b253a72b",
		"char2:78d13a9d1e2210abb524cb48a13577bfbff8bb2d",
		"char2:79607bdb8180d03d702b3c40e4d358137dae9c97",
		"char2:7aa5d950508b2fb5fc3435446ba5320198208e35",
		"char2:7be770941d443664087ad717deaa90fe5a3a4934",
		"char2:7e776df1bac44bc4277079c01b9020df03d1b718",
		"char2:9d94288ecbcde13eec39801aa3c35e89358516ff",
		"char2:9db93903b089f17870dd33c145c1c41c4df95f5a",
		"char2:a4ab1fe0fcee9eae4c9ed60da813a7986cd2aade",
		"char2:bd12c594c1cd028e7cca138504154a6be8e24d75",
		"char2:bfc2022dbba753a5000713772141006a8a4d99b3",
		"char2:d081dcfcdeea33548e63bb231bd18f6aeb06ee24",
		"char2:d36719954499217ed404fa6e77c5e48c057cb8c1",
		"char2:d7abbfba01625bec803788ee42da86461c96e0bd",
		"char2:d96a149a40ef64abfdf4bd71404808f807749654",
		"char2:d9793ddb457316fb3951d1b1092526183270a307",
		"char2:da4f9eacf69feb5eb5bdfad0de1e3c2dc9dee335",
		"char2:dd2f9d4421e67a1373477c950c23b8e4a1d9a376",
		"char2:de2e42ffbb1e4b6097b845829fb477f2e052b41b",
		"char2:df964b8de08309eae4df90663088fdd75cee370e",
		"char2:e2fda8c0ba5a7b3299103fc9d3d57e9831dbc70a",
		"char3:0c7e285035fd37d0cca39c4d21efdf2d69841595",
		"char3:b4ffb1df821da15a0bf03003afb22dce358ac2b5",
		"char3:bcd7b0b10d4c9ef7afea2416c07b545f253e767f",
	}

	markers := make([]*model.FivenetUserLocations, len(userIdentifiers))

	resetMarkers := func() {
		xMin := -3300
		xMax := 4300
		yMin := -3300
		yMax := 5000
		for i := 0; i < len(markers); i++ {
			x := float64(rand.Intn(xMax-xMin+1) + xMin)
			y := float64(rand.Intn(yMax-yMin+1) + yMin)

			job := "ambulance"
			hidden := false
			markers[i] = &model.FivenetUserLocations{
				Identifier: userIdentifiers[i],
				Job:        job,
				Hidden:     &hidden,

				X: &x,
				Y: &y,
			}
		}
	}

	moveMarkers := func() {
		xMin := -100
		xMax := 100
		yMin := -100
		yMax := 100

		for i := 0; i < len(markers); i++ {
			curX := *markers[i].X
			curY := *markers[i].Y

			newX := curX + float64(rand.Intn(xMax-xMin+1)+xMin)
			newY := curY + float64(rand.Intn(yMax-yMin+1)+yMin)

			markers[i].X = &newX
			markers[i].Y = &newY
		}
	}

	resetMarkers()

	counter := 0
	for {
		func() {
			ctx, span := s.tracer.Start(s.ctx, "livemap-gen-users")
			defer span.End()

			if counter >= 60 {
				resetMarkers()
				counter = 0
			} else {
				moveMarkers()
			}

			stmt := tLocs.
				INSERT(
					tLocs.Identifier,
					tLocs.Job,
					tLocs.X,
					tLocs.Y,
					tLocs.Hidden,
				).
				MODELS(markers).
				ON_DUPLICATE_KEY_UPDATE(
					tLocs.X.SET(jet.RawFloat("VALUES(`x`)")),
					tLocs.Y.SET(jet.RawFloat("VALUES(`y`)")),
					tLocs.Hidden.SET(jet.RawBool("VALUES(`hidden`)")),
				)

			_, err := stmt.ExecContext(ctx, s.db)
			if err != nil {
				s.logger.Error("failed to insert/ update random location to locations table", zap.Error(err))
			}

			counter++
			time.Sleep(3 * time.Second)
		}()
	}
}
