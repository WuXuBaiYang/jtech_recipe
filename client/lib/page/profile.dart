import 'package:client/common/logic.dart';
import 'package:flutter/widgets.dart';

/*
* 用户页
* @author wuxubaiyang
* @Time 2022/10/25 9:58
*/
class ProfilePage extends StatefulWidget {
  const ProfilePage({super.key});

  @override
  State<StatefulWidget> createState() => _ProfilePageState();
}

/*
* 用户页-状态
* @author wuxubaiyang
* @Time 2022/10/25 9:58
*/
class _ProfilePageState extends State<ProfilePage> {
  // 逻辑管理
  final _logic = _ProfilePageLogic();

  @override
  Widget build(BuildContext context) {
    return SizedBox();
  }

  @override
  void dispose() {
    _logic.dispose();
    super.dispose();
  }
}

/*
* 用户页-逻辑
* @author wuxubaiyang
* @Time 2022/10/25 9:58
*/
class _ProfilePageLogic extends BaseLogic {}
