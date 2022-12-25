package events

import (
	"net/url"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/yungen-lu/TOC-Project-2022/internal/models"
)

func buildListCarousel(instances []models.Instance) *linebot.CarouselContainer {
	r := make([]*linebot.BubbleContainer, len(instances))
	for i, l := range instances {
		r[i] = buildListFlexMessage(l)
	}
	return &linebot.CarouselContainer{
		Type:     linebot.FlexContainerTypeCarousel,
		Contents: r,
	}
}
func buildListFlexMessage(model models.Instance) *linebot.BubbleContainer {
	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Hero: &linebot.ImageComponent{
			URL:         getImageURL(model.Type),
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
					Text:   model.Name,
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
									Text: model.Type,
									Size: linebot.FlexTextSizeTypeMd,
									Flex: linebot.IntPtr(4),
									Wrap: true,
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
									Text: model.Namespace,
									Size: linebot.FlexTextSizeTypeSm,
									Flex: linebot.IntPtr(4),
									Wrap: true,
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
						Data:  buildQuery(map[string]string{"action": "info", "dbname": model.Name, "ns": model.Namespace}),
					},
				},
				&linebot.ButtonComponent{
					Type:   linebot.FlexComponentTypeButton,
					Style:  linebot.FlexButtonStyleTypeSecondary,
					Height: linebot.FlexButtonHeightTypeSm,
					Action: &linebot.PostbackAction{
						Label: "Delete",
						Data:  buildQuery(map[string]string{"action": "delete", "dbname": model.Name, "ns": model.Namespace}),
					},
				},
			},
		},
	}
}
func buildInfoFlexMessage(model *models.Instance) *linebot.BubbleContainer {
	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   model.Name,
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
									Text: model.Type,
									Size: linebot.FlexTextSizeTypeMd,
									Flex: linebot.IntPtr(4),
									Wrap: true,
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
									Text: model.Namespace,
									Size: linebot.FlexTextSizeTypeSm,
									Flex: linebot.IntPtr(4),
									Wrap: true,
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
									Text: model.User,
									Size: linebot.FlexTextSizeTypeSm,
									Flex: linebot.IntPtr(4),
									Wrap: true,
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
									Text: model.Password,
									Size: linebot.FlexTextSizeTypeSm,
									Flex: linebot.IntPtr(4),
									Wrap: true,
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
									Text: "Endpoint",
									Size: linebot.FlexTextSizeTypeSm,
									Flex: linebot.IntPtr(2),
								},
								&linebot.TextComponent{
									Type: linebot.FlexComponentTypeText,
									Text: model.Endpoint,
									Size: linebot.FlexTextSizeTypeSm,
									Flex: linebot.IntPtr(4),
									Wrap: true,
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
func getImageURL(dbtype string) string {
	switch dbtype {
	case "postgres":
		return "https://i.imgur.com/75ZLw5s.png"
	case "mysql":
		return "https://i.imgur.com/xmUdGYB.png"
	case "redis":
		return "https://i.imgur.com/nJgwQOV.png"
	case "mongodb":
		return "https://i.imgur.com/NLt9Mje.png"
	default:
		return "https://i.imgur.com/75ZLw5s.png"
	}
}
