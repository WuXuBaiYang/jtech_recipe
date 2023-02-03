import 'package:flutter/widgets.dart';

/*
* 逻辑处理基类
* @author wuxubaiyang
* @Time 2022/10/20 10:03
*/
abstract class BaseLogic {
  @mustCallSuper
  void init() {}

  void dispose() {}
}

/*
* 带有逻辑管理结构的状态基类
* @author wuxubaiyang
* @Time 2022/11/2 11:19
*/
abstract class LogicState<T extends StatefulWidget, C extends BaseLogic>
    extends State<T> {
  // 初始化逻辑管理
  C initLogic();

  // 缓存逻辑管理对象
  C? _cacheLogic;

  // 获取逻辑对象
  C get logic => _cacheLogic ??= initLogic();

  @override
  void initState() {
    super.initState();
    logic.init();
  }

  @override
  void dispose() {
    logic.dispose();
    super.dispose();
  }
}
