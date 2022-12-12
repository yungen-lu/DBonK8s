package events

import (
	"net/url"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func buildListCarousel(instances []*Instance) []*linebot.BubbleContainer {
	r := make([]*linebot.BubbleContainer, len(instances))
	for _, l := range instances {
		r = append(r, buildListFlexMessage(l.Name, l.Type, l.Namespace))
	}
	return r
}
func buildListFlexMessage(dbname string, dbtype string, namespace string) *linebot.BubbleContainer {
	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Hero: &linebot.ImageComponent{
			URL:         "https://scdn.line-apps.com/n/channel_devcenter/img/fx/01_1_cafe.png",
			Size:        linebot.FlexImageSizeTypeFull,
			AspectRatio: linebot.FlexImageAspectRatioType1_51to1,
			AspectMode:  linebot.FlexImageAspectModeTypeCover,
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   dbname,
					Size:   linebot.FlexTextSizeTypeXl,
					Weight: linebot.FlexTextWeightTypeBold,
				},
				&linebot.BoxComponent{
					Type:    linebot.FlexComponentTypeBox,
					Layout:  linebot.FlexBoxLayoutTypeVertical,
					Spacing: linebot.FlexComponentSpacingTypeSm,
					Margin:  linebot.FlexComponentMarginTypeLg,
					Contents: []linebot.FlexComponent{
						&linebot.BoxComponent{
							Type:    linebot.FlexComponentTypeBox,
							Layout:  linebot.FlexBoxLayoutTypeBaseline,
							Spacing: linebot.FlexComponentSpacingTypeSm,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type: linebot.FlexComponentTypeText,
									Text: "DBType",
									Size: linebot.FlexTextSizeTypeSm,
									Flex: linebot.IntPtr(2),
								},
								&linebot.TextComponent{
									Type: linebot.FlexComponentTypeText,
									Text: dbtype,
									Size: linebot.FlexTextSizeTypeMd,
									Flex: linebot.IntPtr(4),
								},
							},
						},
						&linebot.BoxComponent{
							Type:    linebot.FlexComponentTypeBox,
							Layout:  linebot.FlexBoxLayoutTypeBaseline,
							Spacing: linebot.FlexComponentSpacingTypeSm,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type: linebot.FlexComponentTypeText,
									Text: "Owner",
									Size: linebot.FlexTextSizeTypeSm,
									Flex: linebot.IntPtr(2),
								},
								&linebot.TextComponent{
									Type: linebot.FlexComponentTypeText,
									Text: namespace,
									Size: linebot.FlexTextSizeTypeSm,
									Flex: linebot.IntPtr(4),
								},
							},
						},
					},
				},
			},
		},
		Footer: &linebot.BoxComponent{
			Type:    linebot.FlexComponentTypeBox,
			Layout:  linebot.FlexBoxLayoutTypeVertical,
			Flex:    linebot.IntPtr(0),
			Spacing: linebot.FlexComponentSpacingTypeSm,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type:   linebot.FlexComponentTypeButton,
					Style:  linebot.FlexButtonStyleTypeSecondary,
					Height: linebot.FlexButtonHeightTypeSm,
					Action: &linebot.PostbackAction{
						Label: "Get Info",
						Data:  buildQuery(map[string]string{"action": "info", "dbname": dbname, "ns": namespace}),
					},
				},
				&linebot.ButtonComponent{
					Type:   linebot.FlexComponentTypeButton,
					Style:  linebot.FlexButtonStyleTypeSecondary,
					Height: linebot.FlexButtonHeightTypeSm,
					Action: &linebot.PostbackAction{
						Label: "Delete",
						Data:  buildQuery(map[string]string{"action": "delete", "dbname": dbname, "ns": namespace}),
					},
				},
			},
		},
	}
}
func buildInfoFlexMessage(dbname string, dbtype string, user string, password string, namespace string) *linebot.BubbleContainer {
	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   dbname,
					Size:   linebot.FlexTextSizeTypeXl,
					Weight: linebot.FlexTextWeightTypeBold,
				},
				&linebot.BoxComponent{
					Type:    linebot.FlexComponentTypeBox,
					Layout:  linebot.FlexBoxLayoutTypeVertical,
					Spacing: linebot.FlexComponentSpacingTypeSm,
					Margin:  linebot.FlexComponentMarginTypeLg,
					Contents: []linebot.FlexComponent{
						&linebot.BoxComponent{
							Type:    linebot.FlexComponentTypeBox,
							Layout:  linebot.FlexBoxLayoutTypeBaseline,
							Spacing: linebot.FlexComponentSpacingTypeSm,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type: linebot.FlexComponentTypeText,
									Text: "DBType",
									Size: linebot.FlexTextSizeTypeSm,
									Flex: linebot.IntPtr(2),
								},
								&linebot.TextComponent{
									Type: linebot.FlexComponentTypeText,
									Text: dbtype,
									Size: linebot.FlexTextSizeTypeMd,
									Flex: linebot.IntPtr(4),
								},
							},
						},
						&linebot.BoxComponent{
							Type:    linebot.FlexComponentTypeBox,
							Layout:  linebot.FlexBoxLayoutTypeBaseline,
							Spacing: linebot.FlexComponentSpacingTypeSm,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type: linebot.FlexComponentTypeText,
									Text: "Owner",
									Size: linebot.FlexTextSizeTypeSm,
									Flex: linebot.IntPtr(2),
								},
								&linebot.TextComponent{
									Type: linebot.FlexComponentTypeText,
									Text: namespace,
									Size: linebot.FlexTextSizeTypeSm,
									Flex: linebot.IntPtr(4),
								},
							},
						},
						&linebot.BoxComponent{
							Type:    linebot.FlexComponentTypeBox,
							Layout:  linebot.FlexBoxLayoutTypeBaseline,
							Spacing: linebot.FlexComponentSpacingTypeSm,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type: linebot.FlexComponentTypeText,
									Text: "User",
									Size: linebot.FlexTextSizeTypeSm,
									Flex: linebot.IntPtr(2),
								},
								&linebot.TextComponent{
									Type: linebot.FlexComponentTypeText,
									Text: user,
									Size: linebot.FlexTextSizeTypeSm,
									Flex: linebot.IntPtr(4),
								},
							},
						},
						&linebot.BoxComponent{
							Type:    linebot.FlexComponentTypeBox,
							Layout:  linebot.FlexBoxLayoutTypeBaseline,
							Spacing: linebot.FlexComponentSpacingTypeSm,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type: linebot.FlexComponentTypeText,
									Text: "Password",
									Size: linebot.FlexTextSizeTypeSm,
									Flex: linebot.IntPtr(2),
								},
								&linebot.TextComponent{
									Type: linebot.FlexComponentTypeText,
									Text: password,
									Size: linebot.FlexTextSizeTypeSm,
									Flex: linebot.IntPtr(4),
								},
							},
						},
					},
				},
			},
		},
	}

}

//	func buildQuery(dbname string, namespace string) string {
//		v := url.Values{}
//		v.Set("dbname", dbname)
//		v.Set("ns", namespace)
//		return v.Encode()
//	}
func buildQuery(m map[string]string) string {
	v := url.Values{}
	for key, value := range m {
		v.Set(key, value)
	}
	return v.Encode()

}
