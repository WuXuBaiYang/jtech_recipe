import 'package:client/manage/event.dart';
import 'package:flutter/material.dart';

/*
* 全局样式控制事件
* @author wuxubaiyang
* @Time 2022/4/1 15:14
*/
class ThemeEvent extends EventModel {
  // 全局样式
  final ThemeData themeData;

  const ThemeEvent({required this.themeData});
}
