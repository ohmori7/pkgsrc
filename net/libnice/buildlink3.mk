# $NetBSD: buildlink3.mk,v 1.32 2018/10/24 21:11:47 leot Exp $

BUILDLINK_TREE+=	libnice

.if !defined(LIBNICE_BUILDLINK3_MK)
LIBNICE_BUILDLINK3_MK:=

BUILDLINK_API_DEPENDS.libnice+=	libnice>=0.0.9
BUILDLINK_ABI_DEPENDS.libnice+=	libnice>=0.1.14nb7
BUILDLINK_PKGSRCDIR.libnice?=	../../net/libnice

.include "../../devel/glib2/buildlink3.mk"
.include "../../security/gnutls/buildlink3.mk"
.include "../../net/gupnp-igd/buildlink3.mk"
.endif	# LIBNICE_BUILDLINK3_MK

BUILDLINK_TREE+=	-libnice
