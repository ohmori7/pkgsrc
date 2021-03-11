# $NetBSD: buildlink3.mk,v 1.34 2021/01/01 08:24:51 ryoon Exp $

BUILDLINK_TREE+=	sword

.if !defined(SWORD_BUILDLINK3_MK)
SWORD_BUILDLINK3_MK:=

BUILDLINK_API_DEPENDS.sword+=	sword>=1.8.1
BUILDLINK_ABI_DEPENDS.sword?=	sword>=1.8.1nb9
BUILDLINK_PKGSRCDIR.sword?=	../../misc/sword

.include "../../textproc/icu/buildlink3.mk"
.include "../../textproc/libclucene/buildlink3.mk"
.include "../../devel/zlib/buildlink3.mk"
.include "../../www/curl/buildlink3.mk"
.endif	# SWORD_BUILDLINK3_MK

BUILDLINK_TREE+=	-sword
