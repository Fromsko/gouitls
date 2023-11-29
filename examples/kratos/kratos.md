## 用户管理

`app/manage/user`

```shell
kratos new app/manage/user --nomod -r https://gitee.com/go-kratos/kratos-layout
```

```shell
kratos proto server api/manage/user_service/user_service.proto -t app/manage/user/internal/service
```

---

```shell
import "google/api/annotations.proto";
```

```shell
# 增加 proto 文件
kratos proto add api/manage/exam_service/exam.proto
kratos proto add api/manage/question_bank_service/question_bank.proto
kratos proto add api/manage/grade_service/grade.proto
kratos proto add api/manage/class_service/class.proto
kratos proto add api/manage/student_service/student.proto
kratos proto add api/manage/course_service/course.proto
kratos proto add api/manage/teacher_service/teacher.proto

kratos new usermini/manage --nomod https://gitee.com/huoyingwhw/kratos-gin

# 大仓
kratos new app/manage/exam --nomod -r https://gitee.com/go-kratos/kratos-layout
kratos new app/manage/question_bank --nomod -r https://gitee.com/go-kratos/kratos-layout
kratos new app/manage/grade --nomod -r https://gitee.com/go-kratos/kratos-layout
kratos new app/manage/class --nomod -r https://gitee.com/go-kratos/kratos-layout
kratos new app/manage/student --nomod -r https://gitee.com/go-kratos/kratos-layout
kratos new app/manage/course --nomod -r https://gitee.com/go-kratos/kratos-layout
kratos new app/manage/teacher --nomod -r https://gitee.com/go-kratos/kratos-layout

# 服务端
kratos proto server api/manage/exam_service/exam.proto
kratos proto server api/manage/question_bank_service/question_bank.proto
kratos proto server api/manage/grade_service/grade.proto
kratos proto server api/manage/class_service/class.proto
kratos proto server api/manage/student_service/student.proto
kratos proto server api/manage/course_service/course.proto
kratos proto server api/manage/teacher_service/teacher.proto

# 客户端
kratos proto client api/manage/exam_service/exam.proto
kratos proto client api/manage/question_bank_service/question_bank.proto
kratos proto client api/manage/grade_service/grade.proto
kratos proto client api/manage/class_service/class.proto
kratos proto client api/manage/student_service/student.proto
kratos proto client api/manage/course_service/course.proto
kratos proto client api/manage/teacher_service/teacher.proto
```
