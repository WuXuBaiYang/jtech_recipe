import 'package:client/common/model.dart';
import 'package:client/tool/date.dart';

import 'model.dart';

/*
* 用户信息
* @author wuxubaiyang
* @Time 2022/9/12 19:10
*/
class UserModel extends BaseModel with BaseInfo {
  // im用户id
  final String imUserId;

  // im授权信息
  final String imToken;

  // im过期时间
  final DateTime imExpired;

  // 用户名
  final String userName;

  // 用户信息
  final UserProfileModel? profile;

  UserModel.from(obj)
      : imUserId = obj?["imUserId"] ?? "",
        imToken = obj?["imToken"] ?? "",
        imExpired = DateTime.fromMillisecondsSinceEpoch(
            (obj?["imExpired"] ?? 0) * 1000),
        userName = obj?["userName"] ?? "",
        profile = null != obj?["profile"]
            ? UserProfileModel.from(obj?["profile"] ?? {})
            : null {
    initialBaseInfo(obj);
  }

  @override
  Map<String, dynamic> to() => {
        ...baseInfoMap,
        "imUserId": imUserId,
        "imToken": imToken,
        "imExpired": imExpired.millisecondsSinceEpoch ~/ 1000,
        "userName": userName,
      };
}

/*
* 用户信息
* @author wuxubaiyang
* @Time 2022/9/12 19:38
*/
class UserProfileModel extends BaseModel with BaseInfo, CreatorInfo {
  // 昵称
  final String nickName;

  // 头像
  final String avatar;

  // 手机号
  final String telephone;

  // 简介
  final String bio;

  // 地址
  final String address;

  // 经纬度（逗号分隔）
  final String location;

  // 职业
  final String profession;

  // 邮箱
  final String email;

  // 性别（0未选择|1男|2女|3武装直升机）
  final GenderType gender;

  // 生日
  final DateTime birthday;

  UserProfileModel.from(obj)
      : nickName = obj?["nickName"] ?? "",
        avatar = obj?["avatar"] ?? "",
        telephone = obj?["telephone"] ?? "",
        bio = obj?["bio"] ?? "",
        address = obj?["address"] ?? "",
        location = obj?["location"] ?? "",
        profession = obj?["profession"] ?? "",
        email = obj?["email"] ?? "",
        gender = GenderType.values[obj?["gender"] ?? 0],
        birthday = DateTool.parseDate(obj?["birthday"] ?? "") ?? DateTime(0) {
    initialBaseInfo(obj);
    initialCreatorInfo(obj);
  }

  @override
  Map<String, dynamic> to() => {
        "nickName": nickName,
        "avatar": avatar,
        "telephone": telephone,
        "bio": bio,
        "address": address,
        "location": location,
        "profession": profession,
        "email": email,
        "gender": gender.index,
        "birthday": DateTool.formatDate(DatePattern.fullDate, birthday),
      };
}

/*
* 用户性别枚举
* @author wuxubaiyang
* @Time 2022/9/12 20:20
*/
enum GenderType {
  // 未选择
  unknown,
  // 男性
  male,
  // 女性
  female,
  // 武装直升机
  militaryHelicopter,
}

/*
* 用户性别枚举扩展
* @author wuxubaiyang
* @Time 2022/9/12 20:24
*/
extension GenderTypeExtension on GenderType {
  // 获取性别名称
  String get name {
    switch (this) {
      case GenderType.unknown:
        return "未选择";
      case GenderType.male:
        return "男";
      case GenderType.female:
        return "女";
      case GenderType.militaryHelicopter:
        return "武装直升机";
    }
  }
}
