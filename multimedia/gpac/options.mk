# $NetBSD: options.mk,v 1.3 2020/09/19 03:38:07 schmonz Exp $

PKG_OPTIONS_VAR=	PKG_OPTIONS.gpac
PKG_SUPPORTED_OPTIONS=	x11
PKG_SUGGESTED_OPTIONS=	x11

.include "../../mk/bsd.options.mk"

PLIST_VARS+=		x11
.if !empty(PKG_OPTIONS:Mx11)
CONFIGURE_ARGS+=	--X11-path=${X11BASE}
PLIST.x11=		yes
.include "../../graphics/glu/buildlink3.mk"
.include "../../graphics/MesaLib/buildlink3.mk"
.include "../../x11/libX11/buildlink3.mk"
.include "../../x11/libXv/buildlink3.mk"
.endif
