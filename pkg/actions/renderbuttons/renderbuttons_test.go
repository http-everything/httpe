package renderbuttons_test

import (
	"fmt"
	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/actions/renderbuttons"
	"http-everything/httpe/pkg/requestdata"
	"http-everything/httpe/pkg/rules"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func TestRenderButtons(t *testing.T) {
	cases := []struct {
		name string
		rule rules.Rule
	}{
		{
			name: "Embedded Template",
			rule: rules.Rule{
				Name: "Render My Buttons",
				Do: &rules.Do{
					RenderButtons: []rules.Button{
						{
							Name:    "Button 1",
							URL:     "/action/1",
							Classes: "btn-lg",
						},
						{
							Name: "Button 2",
							URL:  "/action/2",
						},
					},
				},
			},
		},
		{
			name: "Template from file system",
			rule: rules.Rule{
				Name: "Render My Buttons",
				Do: &rules.Do{
					RenderButtons: []rules.Button{
						{
							Name:    "Button 1",
							URL:     "/action/1",
							Classes: "btn-lg",
						},
						{
							Name: "Button 2",
							URL:  "/action/2",
						},
					},
					Args: rules.Args{
						Template: "./buttons.tpl.html",
					},
				},
			},
		},
	}
	var actioner actions.Actioner = renderbuttons.RenderButtons{}
	for _, tc := range cases {
		t.Logf("%s: %s\n", tc.name, tc.rule.Name)

		t.Run(tc.name, func(t *testing.T) {
			actionResp, err := actioner.Execute(tc.rule, requestdata.Data{})
			require.NoError(t, err)

			// Assert html is valid and title is correct
			doc, err := html.Parse(strings.NewReader(actionResp.SuccessBody))
			assert.NoError(t, err, "HTML parsing failed")
			assert.Contains(t, actionResp.SuccessBody, "<title>Render My Buttons</title>")

			// Find all buttons inside the HTML
			type button struct {
				Name      string
				Attribute string
			}
			var buttons []button
			var btnFn func(*html.Node)
			btnFn = func(n *html.Node) {
				if n.Type == html.ElementNode && n.Data == "button" {
					for _, b := range n.Attr {
						if b.Key == "@click" && strings.Contains(b.Val, "fetchUrl") {
							//adds a new button entry when the class starts with "btn"
							buttons = append(buttons, button{
								Name:      strings.TrimSuffix(n.FirstChild.Data, "\n"),
								Attribute: b.Val,
							})
						}
					}
				}

				// traverses the HTML of the webpage from the first child node
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					btnFn(c)
				}
			}
			btnFn(doc)

			// loops through the links slice
			for _, v := range buttons {
				t.Logf("Button name: '%s', attr: '%s'\n", v.Name, v.Attribute)
			}
			for i, c := range tc.rule.Do.RenderButtons {
				assert.Equal(t, c.Name, buttons[i].Name)
				assert.Contains(t, buttons[i].Attribute, fmt.Sprintf("fetchUrl('%s')", c.URL))
			}
		})
	}
}
