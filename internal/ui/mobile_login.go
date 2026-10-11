package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"

	"gioui.org/widget/material"
)

func (u *UI) layoutMobileLoginCard(gtx C) D {
	p := u.pal
	th := u.th

	return layout.Background{}.Layout(gtx,
		func(gtx C) D {
			return fill(gtx, p.Frame)
		},
		func(gtx C) D {
			return layout.Inset{
				Top:    unit.Dp(24),
				Bottom: unit.Dp(24),
				Left:   unit.Dp(20),
				Right:  unit.Dp(20),
			}.Layout(gtx, func(gtx C) D {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return material.H6(
							th,
							"bubuuuuu",
						).Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
					layout.Rigid(func(gtx C) D {
						return material.Body2(
							th,
							"phone number",
						).Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
					layout.Rigid(func(gtx C) D {
						return material.Editor(
							th,
							&u.phoneEditor,
							"hihihihihihi",
						).Layout(gtx)
					}),
				)
			})
		},
	)
}
func (u *UI) layoutMobileLogin(gtx C) D {
	p := u.pal

	dims := fill(gtx, p.Frame)

	gtx.Constraints.Min = gtx.Constraints.Max

	layout.Center.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Max.X = min(
			gtx.Constraints.Max.X-gtx.Dp(32),
			gtx.Dp(480),
		)

		return layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Rigid(u.layoutMobileLoginCard),
			layout.Rigid(layout.Spacer{Height: 24}.Layout),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return drawIcon(
							gtx,
							icLock,
							14,
							p.TextSecondary,
						)
					}),
					layout.Rigid(layout.Spacer{Width: 5}.Layout),
					layout.Rigid(
						u.label(
							13,
							u.locale.Text("PolicyE2E"),
							p.TextSecondary,
						).Layout,
					),
				)
			}),
		)
	})

	return dims
}