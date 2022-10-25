import 'package:client/common/common.dart';
import 'package:client/common/logic.dart';
import 'package:client/manage/auth.dart';
import 'package:client/manage/router.dart';
import 'package:client/widget/avatar.dart';
import 'package:client/widget/oss.dart';
import 'package:flutter/material.dart';

/*
* 首页
* @author wuxubaiyang
* @Time 2022/9/8 15:01
*/
class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<StatefulWidget> createState() => _HomePageState();
}

/*
* 首页-状态
* @author wuxubaiyang
* @Time 2022/9/8 15:02
*/
class _HomePageState extends State<HomePage> {
  // 逻辑管理
  final _logic = _HomeLogic();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('每日定食'),
        actions: [
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 8),
            child: OSSAvatar(
              object: authManage.userInfo.avatar,
              avatarSize: AvatarSize.small,
              onTap: () => routerManage.pushNamed(RoutePath.profile),
            ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {},
      ),
    );
  }

  @override
  void dispose() {
    _logic.dispose();
    super.dispose();
  }
}

/*
* 首页-逻辑
* @author wuxubaiyang
* @Time 2022/10/25 9:18
*/
class _HomeLogic extends BaseLogic {}
