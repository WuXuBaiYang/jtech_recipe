import 'package:client/common/logic.dart';
import 'package:client/manage/auth.dart';
import 'package:client/widget/avatar.dart';
import 'package:client/widget/oss.dart';
import 'package:flutter/material.dart';

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
    return Scaffold(
      body: CustomScrollView(
        physics: const BouncingScrollPhysics(),
        slivers: [
          SliverAppBar(
            scrolledUnderElevation: 0,
            expandedHeight: 130,
            pinned: true,
            flexibleSpace: FlexibleSpaceBar(
              titlePadding: const EdgeInsets.symmetric(horizontal: 35),
              title: _buildProfile(),
              expandedTitleScale: 1.3,
            ),
          ),
          SliverList(
            delegate: SliverChildListDelegate(
              [
                _buildUserProfile(),
                _buildSettings(),
              ],
            ),
          ),
        ],
      ),
    );
  }

  // 构建基础用户信息
  Widget _buildProfile() {
    final user = authManage.userInfo;
    var bio = user.bio;
    if (bio.isEmpty) bio = "可能需要写点什么";
    return ListTile(
      horizontalTitleGap: 0,
      leading: OSSAvatar(
        object: user.avatar,
        avatarSize: AvatarSize.small,
      ),
      title: Text(user.nickName),
      subtitle: Text(
        bio,
        style: const TextStyle(fontSize: 10, color: Colors.black54),
        overflow: TextOverflow.ellipsis,
        softWrap: true,
        maxLines: 1,
      ),
    );
  }

  // 构建用户信息
  Widget _buildUserProfile() {
    return Column(
      children: [],
    );
  }

  // 构建系统设置
  Widget _buildSettings() {
    return Text('a');
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
