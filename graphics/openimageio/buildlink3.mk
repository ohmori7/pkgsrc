# $NetBSD: buildlink3.mk,v 1.51 2021/01/19 16:02:25 nia Exp $

BUILDLINK_TREE+=	openimageio

.if !defined(OPENIMAGEIO_BUILDLINK3_MK)
OPENIMAGEIO_BUILDLINK3_MK:=

BUILDLINK_API_DEPENDS.openimageio+=	openimageio>=1.8.15
BUILDLINK_ABI_DEPENDS.openimageio?=	openimageio>=1.8.15nb21
BUILDLINK_PKGSRCDIR.openimageio?=	../../graphics/openimageio

.include "../../devel/boost-libs/buildlink3.mk"
.include "../../graphics/freetype2/buildlink3.mk"
.include "../../graphics/glew/buildlink3.mk"
.include "../../graphics/ilmbase/buildlink3.mk"
.include "../../graphics/libwebp/buildlink3.mk"
.include "../../graphics/openjpeg/buildlink3.mk"
.include "../../graphics/openexr/buildlink3.mk"
.include "../../graphics/tiff/buildlink3.mk"
.endif	# OPENIMAGEIO_BUILDLINK3_MK

BUILDLINK_TREE+=	-openimageio
