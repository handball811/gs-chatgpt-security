# 入力と応答の記録


## User
### パスワードに利用される技術
https://auth0.com/blog/jp-hashing-passwords-one-way-road-to-security/
- Hash化することが前提
- ハッシュから元の文字列を生成することができない
- 一般的なハッシュアルゴリズム
  - md5, sha256
  - 

### UserCreate struct

#### 1
```入力
以下のUserCreateをUserCreateOpのinterfaceに従う形で作成してください。
Createメソッドは何らかのデータベースへの保存機能を定義してください。
package userdraft

import "context"

type User struct {
	Email       string
	DisplayName string
}

type UserCreateOp interface {
	SetEmail(email string) UserCreateOp
	SetPassword(password string) UserCreateOp
	SetDisplayName(password string) UserCreateOp
	Create(context.Context) (*User, error)
}

type UserCreate struct {}
```

```
package userdraft

import "context"

type User struct {
	Email       string
	DisplayName string
}

type UserCreateOp interface {
	SetEmail(email string) UserCreateOp
	SetPassword(password string) UserCreateOp
	SetDisplayName(password string) UserCreateOp
	Create(context.Context) (*User, error)
}

type UserCreate struct {
	email       string
	password    string
	displayName string
}

func (uc *UserCreate) SetEmail(email string) UserCreateOp {
	uc.email = email
	return uc
}

func (uc *UserCreate) SetPassword(password string) UserCreateOp {
	uc.password = password
	return uc
}

func (uc *UserCreate) SetDisplayName(displayName string) UserCreateOp {
	uc.displayName = displayName
	return uc
}

func (uc *UserCreate) Create(ctx context.Context) (*User, error) {
	// Save user to database
	return &User{
		Email:       uc.email,
		DisplayName: uc.displayName,
	}, nil
}
```

#### 2

```
データベースに保管するためのコードを追加してください
```

```
package userdraft

import (
	"context"
	"database/sql"
)

type User struct {
	Email       string
	DisplayName string
}

type UserCreateOp interface {
	SetEmail(email string) UserCreateOp
	SetPassword(password string) UserCreateOp
	SetDisplayName(password string) UserCreateOp
	Create(context.Context) (*User, error)
}

type UserCreate struct {
	db          *sql.DB
	email       string
	password    string
	displayName string
}

func NewUserCreate(db *sql.DB) UserCreateOp {
	return &UserCreate{
		db: db,
	}
}

func (uc *UserCreate) SetEmail(email string) UserCreateOp {
	uc.email = email
	return uc
}

func (uc *UserCreate) SetPassword(password string) UserCreateOp {
	uc.password = password
	return uc
}

func (uc *UserCreate) SetDisplayName(displayName string) UserCreateOp {
	uc.displayName = displayName
	return uc
}

func (uc *UserCreate) Create(ctx context.Context) (*User, error) {
	// Save user to database
	query := "INSERT INTO users (email, password, display_name) VALUES (?, ?, ?)"
	result, err := uc.db.ExecContext(ctx, query, uc.email, uc.password, uc.displayName)
	if err != nil {
		return nil, err
	}

	// Get the newly inserted user ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &User{
		Email:       uc.email,
		DisplayName: uc.displayName,
	}, nil
}
```
#### 補足
idはしようされていないためエラーとなる

#### 3

```
パスワードがそのまま保管されてしまっているので閲覧できないようにしてください
```

```
package userdraft

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email       string
	DisplayName string
}

type UserCreateOp interface {
	SetEmail(email string) UserCreateOp
	SetPassword(password string) UserCreateOp
	SetDisplayName(password string) UserCreateOp
	Create(context.Context) (*User, error)
}

type UserCreate struct {
	db          *sql.DB
	email       string
	password    string
	displayName string
}

func NewUserCreate(db *sql.DB) UserCreateOp {
	return &UserCreate{
		db: db,
	}
}

func (uc *UserCreate) SetEmail(email string) UserCreateOp {
	uc.email = email
	return uc
}

func (uc *UserCreate) SetPassword(password string) UserCreateOp {
	// Hash the password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uc
	}
	uc.password = string(hashed)
	return uc
}

func (uc *UserCreate) SetDisplayName(displayName string) UserCreateOp {
	uc.displayName = displayName
	return uc
}

func (uc *UserCreate) Create(ctx context.Context) (*User, error) {
	// Save user to database
	query := "INSERT INTO users (email, password, display_name) VALUES (?, ?, ?)"
	result, err := uc.db.ExecContext(ctx, query, uc.email, uc.password, uc.displayName)
	if err != nil {
		return nil, err
	}

	// Get the newly inserted user ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &User{
		Email:       uc.email,
		DisplayName: uc.displayName,
	}, nil
}

```
#### 補足
https://pkg.go.dev/golang.org/x/crypto/bcrypt


#### 4
```
データベースへの保管はUsetModelRepository interfaceを定義して保管してください
```
```
package userdraft

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email       string
	DisplayName string
}

type UserModelRepository interface {
	Create(ctx context.Context, email, password, displayName string) (*User, error)
}

type UserCreateOp interface {
	SetEmail(email string) UserCreateOp
	SetPassword(password string) UserCreateOp
	SetDisplayName(password string) UserCreateOp
	Create(context.Context) (*User, error)
}

type UserCreate struct {
	userModelRepository UserModelRepository
	email               string
	password            string
	displayName         string
}

func NewUserCreate(userModelRepository UserModelRepository) UserCreateOp {
	return &UserCreate{
		userModelRepository: userModelRepository,
	}
}

func (uc *UserCreate) SetEmail(email string) UserCreateOp {
	uc.email = email
	return uc
}

func (uc *UserCreate) SetPassword(password string) UserCreateOp {
	// Hash the password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uc
	}
	uc.password = string(hashed)
	return uc
}

func (uc *UserCreate) SetDisplayName(displayName string) UserCreateOp {
	uc.displayName = displayName
	return uc
}

func (uc *UserCreate) Create(ctx context.Context) (*User, error) {
	// Save user to database
	return uc.userModelRepository.Create(ctx, uc.email, uc.password, uc.displayName)
}
```

#### 5
```
Userのemailとpasswordの組み合わせを検証するUserCheckOp interfaceを作成して下さい。その際UserModelRepositoryにGetByEmailメソッドを追加してください。
UserCheckOpとUserModelRepositoryのみを出力してください。
```
```
package userdraft

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email       string
	DisplayName string
	Password    string
}

type UserModelRepository interface {
	Create(ctx context.Context, email, password, displayName string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type UserCheckOp interface {
	SetEmail(email string) UserCheckOp
	SetPassword(password string) UserCheckOp
	Check(ctx context.Context) (*User, error)
}

type UserCheck struct {
	userModelRepository UserModelRepository
	email               string
	password            string
}

func NewUserCheck(userModelRepository UserModelRepository) UserCheckOp {
	return &UserCheck{
		userModelRepository: userModelRepository,
	}
}

func (uc *UserCheck) SetEmail(email string) UserCheckOp {
	uc.email = email
	return uc
}

func (uc *UserCheck) SetPassword(password string) UserCheckOp {
	uc.password = password
	return uc
}

func (uc *UserCheck) Check(ctx context.Context) (*User, error) {
	// Get user from database
	user, err := uc.userModelRepository.GetByEmail(ctx, uc.email)
	if err != nil {
		return nil, err
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(uc.password))
	if err != nil {
		return nil, err
	}

	return user, nil
}
```

```
次のように、上記のコードを使用して、データベースにユーザーを保存する処理を実装することができます。

import (
	"context"
	"os"

	"github.com/handball811/gs-chatgpt-security/internal"
)

func main() {
	// 環境変数からAPIキーとプロジェクトIDを取得する
	apikey := os.Getenv("API_KEY")
	projectid := os.Getenv("PROJECT_ID")

	auth, err := internal.NewAuth(apikey, projectid)
	if err != nil {
		// エラー処理を行う
	}

	router := gin.Default()

	router.POST("/signup", func(c *gin.Context) {
		var form SignupForm
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ユーザーをデータベースに保存する
		authParam, err := auth.Signup(context.Background(), &internal.AuthCreateParam{
			Email:       form.Username,
			Password:    form.Password,
			DisplayName: form.Username,
		})
		if err != nil {
			// エラー処理を行う
		}

		c.JSON(http.StatusOK, gin.H{"message": "Successfully signed up!", "user": authParam})
	})

	router.Run()
}
上記のコードでは、NewAuth関数を呼び出す際に、環境変数からAPIキーとプロジェクトIDを取得しています。また、Signup関数を呼び出す際には、AuthCreateParam構造体を渡しています。この構造体には、ユーザーが入力したサインアップフォームの内容が格納されています。

Signup関数を呼び出すことで、データベースに新しいユーザーを保存することができます。ユーザーを保存できた場合は、AuthParam構造体を返すようになっています。この構造体には、新しく作成されたユーザーのEmailアドレスや表示名
```