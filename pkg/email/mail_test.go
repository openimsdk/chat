package email

//func TestEmail(T *testing.T) {
//	if err := InitConfig(); err != nil {
//		fmt.Fprintf(os.Stderr, "\n\nexit -1: \n%+v\n\n", err)
//		os.Exit(-1)
//	}
//	tests := []struct {
//		name string
//		ctx  context.Context
//		mail string
//		code string
//		want error
//	}{
//		{
//			name: "success send email",
//			ctx:  context.Background(),
//			mail: "test@gmail.com",
//			code: "5555",
//			want: errors.New("nil"),
//		},
//		{
//			name: "fail send email",
//			ctx:  context.Background(),
//			mail: "",
//			code: "5555",
//			want: errors.New("dial tcp :0: connectex: The requested address is not valid in its context."),
//		},
//	}
//	mail := NewMail()
//
//	for _, tt := range tests {
//		T.Run(tt.name, func(t *testing.T) {
//			if got := mail.SendMail(tt.ctx, tt.mail, tt.code); errors.Is(got, tt.want) {
//				t.Errorf("%v have a err,%v", tt.name, tt.want)
//			}
//		})
//	}
//}
//
//func InitConfig() error {
//	yam, err := ioutil.ReadFile("../../config/config.yaml")
//	if err != nil {
//		return err
//	}
//	err = yaml.Unmarshal(yam, &config.Config)
//	if err != nil {
//		return err
//	}
//	return nil
//}
