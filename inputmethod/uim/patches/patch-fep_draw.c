* temporal fix for drawing preedit region in FEP.

*** fep/draw.c.orig	Thu May  2 19:38:00 2013
--- fep/draw.c	Thu May  2 20:13:04 2013
***************
*** 694,700 ****
    if ((g_opt.no_report_cursor && preedit->width == prev_preedit->width) || (!g_opt.no_report_cursor && is_eq_region())) {
      int eq_width_rev = compare_preedit_rev(preedit, prev_preedit);
      debug2(("eq_width_rev = %d\n", eq_width_rev));
!     draw_subpreedit(preedit, eq_width, preedit->width - eq_width_rev);
      if (g_opt.no_report_cursor) {
        put_move_cur(preedit->width - eq_width_rev, preedit->cursor);
      } else {
--- 694,700 ----
    if ((g_opt.no_report_cursor && preedit->width == prev_preedit->width) || (!g_opt.no_report_cursor && is_eq_region())) {
      int eq_width_rev = compare_preedit_rev(preedit, prev_preedit);
      debug2(("eq_width_rev = %d\n", eq_width_rev));
!     draw_subpreedit(preedit, eq_width, preedit->width);
      if (g_opt.no_report_cursor) {
        put_move_cur(preedit->width - eq_width_rev, preedit->cursor);
      } else {

