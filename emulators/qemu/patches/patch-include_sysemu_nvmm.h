$NetBSD: patch-include_sysemu_nvmm.h,v 1.3 2021/03/06 11:19:34 reinoud Exp $

--- include/sysemu/nvmm.h.orig	2021-03-05 22:29:22.991663471 +0000
+++ include/sysemu/nvmm.h
@@ -0,0 +1,26 @@
+/*
+ * Copyright (c) 2018-2019 Maxime Villard, All rights reserved.
+ *
+ * NetBSD Virtual Machine Monitor (NVMM) accelerator support.
+ *
+ * This work is licensed under the terms of the GNU GPL, version 2 or later.
+ * See the COPYING file in the top-level directory.
+ */
+
+#ifndef QEMU_NVMM_H
+#define QEMU_NVMM_H
+
+#include "config-host.h"
+#include "qemu-common.h"
+
+#ifdef CONFIG_NVMM
+
+int nvmm_enabled(void);
+
+#else /* CONFIG_NVMM */
+
+#define nvmm_enabled() (0)
+
+#endif /* CONFIG_NVMM */
+
+#endif /* CONFIG_NVMM */
