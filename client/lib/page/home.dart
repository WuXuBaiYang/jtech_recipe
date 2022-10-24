import 'package:client/common/common.dart';
import 'package:client/manage/auth.dart';
import 'package:client/manage/router.dart';
import 'package:client/manage/tag.dart';
import 'package:client/manage/theme.dart';
import 'package:client/widget/avatar.dart';
import 'package:client/widget/oss.dart';
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';

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
  // 根节点key
  final scaffoldKey = GlobalKey<ScaffoldState>();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      key: scaffoldKey,
      appBar: AppBar(
        leading: IconButton(
          icon: const FaIcon(FontAwesomeIcons.barsStaggered),
          onPressed: () => scaffoldKey.currentState?.openDrawer(),
        ),
        title: const Text('每日定食'),
      ),
      drawerEnableOpenDragGesture: false,
      drawer: _buildDrawerMenu(context),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          if (themeManage.currentTheme.brightness == Brightness.light) {
            themeManage.switchTheme(ThemeType.dark);
          } else {
            themeManage.switchTheme(ThemeType.light);
          }
          setState(() {});
        },
      ),
    );
  }

  // 构建侧滑菜单
  Widget _buildDrawerMenu(BuildContext context) {
    return Drawer(
      child: ListView(
        padding: EdgeInsets.zero,
        children: [
          _buildDrawerHeader(context),
          ListTile(
            leading: const FaIcon(FontAwesomeIcons.message),
            title: const Text('消息通知'),
            onTap: () => routerManage.pushNamed(RoutePath.notice),
          ),
        ],
      ),
    );
  }

  // 构建侧滑菜单头
  Widget _buildDrawerHeader(BuildContext context) {
    final userInfo = authManage.userInfo;
    return DrawerHeader(
      decoration: BoxDecoration(
        color: themeManage.primaryColor,
      ),
      padding: EdgeInsets.zero,
      child: Center(
        child: ListTile(
          contentPadding: EdgeInsets.zero,
          textColor: Colors.white,
          horizontalTitleGap: 0,
          leading: OSSAvatar(
            object: userInfo.avatar,
            avatarSize: AvatarSize.large,
          ),
          title: Text(userInfo.nickName),
          subtitle: Text(userInfo.bio),
          onTap: () => routerManage.pushNamed(RoutePath.myProfile),
        ),
      ),
    );
  }
}
