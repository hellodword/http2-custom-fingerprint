# http2-custom-fingerprint

> modify `golang.org/x/net/http2` fingerprint for mimicry purposes

## layout

- `patches` directory contains patches for each tag
- `.github/workflows/patch.yml` patch all tags automatically

## usage

1. replace `golang.org/x/net` with `github.com/hellodword/http2-custom-fingerprint` in `go.mod`, you can find the generated versions in [`revision.txt`](./revision.txt)
2. use custom functions in `http2.Transport`:

```go
http2.Transport{
	CustomInitialTransportConnFlow: func(u uint32) uint32 {
		fmt.Println("original InitialTransportConnFlow", u)
		return 15663105
	},
	CustomInitialSettings: func(s []http2.Setting) []http2.Setting {
		fmt.Println("original InitialSettings", s)
		return []http2.Setting{
			{ID: http2.SettingHeaderTableSize, Val: 65536},
			{ID: http2.SettingEnablePush, Val: 0},
			{ID: http2.SettingInitialWindowSize, Val: 6291456},
			{ID: http2.SettingMaxHeaderListSize, Val: 262144},
		}
	},
	CustomFirstHeadersFrameParam: func(hfp http2.HeadersFrameParam) http2.HeadersFrameParam {
		fmt.Println("original FirstWriteHeaders", hfp)
		hfp.Priority = http2.PriorityParam{
			Weight:    255,
			StreamDep: 0,
			Exclusive: true,
		}
		return hfp
	},
	CustomHeaders: func(s [][2]string) [][2]string {
		fmt.Println("original Headers", s)
        // sort headers as you want
		return s
	},
}
```
