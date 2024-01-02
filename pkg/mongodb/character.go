// Author liyan
// Date 2020-03-27 5:24 下午
// Mail liyana@hualala.com
// org base cloud platform

package mongodb

//
//type strCase bool
//
//const (
//	//lower strCase = false
//	upper strCase = true
//)

//var commonInitialisms = []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "TTL", "UI", "UID", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
//var commonInitialismsReplacer *strings.Replacer

//func toSnake(name string) string {
//	var (
//		value = commonInitialismsReplacer.Replace(name)
//		buf   = bytes.NewBufferString("")
//
//		lastCase, currCase, nextCase strCase
//	)
//
//	for i, v := range value[:len(value)-1] {
//		nextCase = strCase(value[i+1] >= 'A' && value[i+1] <= 'Z')
//		if i > 0 {
//			if currCase == upper {
//				if lastCase == upper && nextCase == upper {
//					buf.WriteRune(v)
//				} else {
//					if value[i-1] != '_' && value[i+1] != '_' {
//						buf.WriteRune('_')
//					}
//					buf.WriteRune(v)
//				}
//			} else {
//				buf.WriteRune(v)
//			}
//		} else {
//			currCase = upper
//			buf.WriteRune(v)
//		}
//		lastCase = currCase
//		currCase = nextCase
//	}
//
//	buf.WriteByte(value[len(value)-1])
//
//	s := strings.ToLower(buf.String())
//	return s
//}

//func init() {
//	var commonInitialismsForReplacer []string
//	for _, initialism := range commonInitialisms {
//		commonInitialismsForReplacer = append(commonInitialismsForReplacer, initialism, strings.Title(strings.ToLower(initialism)))
//	}
//	commonInitialismsReplacer = strings.NewReplacer(commonInitialismsForReplacer...)
//}
