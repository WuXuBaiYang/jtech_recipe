import 'package:client/manage/auth.dart';
import 'package:client/manage/oss.dart';
import 'package:client/manage/theme.dart';
import 'package:client/model/user.dart';
import 'package:client/tool/tool.dart';
import 'package:client/widget/avatar.dart';
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
  late UserModel user;

  @override
  void initState() {
    super.initState();
    user = authManage.userInfo;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('首页'),
      ),
      body: Center(
        child: Avatar(
          user: user,
        ),
      ),
      floatingActionButton: FloatingActionButton(
        child: const Icon(Icons.device_hub),
        onPressed: () async {
          setState(() {
            user = UserModel.from({'avatar': 'test_avatar.jpg'});
          });
          themeManage.switchTheme(ThemeType.dark);
        },
      ),
    );
  }
}
