// socks5
package utils

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

//SOCKS version
const SOCKS5 = 5

// SOCKS request commands as defined in RFC 1928 section 4.
const (
	CmdConnect      = 1
	CmdBind         = 2
	CmdUDPAssociate = 3
)

// SOCKS address types as defined in RFC 1928 section 5.
const (
	AtypIPv4       = 1
	AtypDomainName = 3
	AtypIPv6       = 4
)

func Socks5HandShake(r *bufio.Reader, conn net.Conn) error {
	version, _ := r.ReadByte()
	if version != SOCKS5 {
		return errors.New("not SOCKS5 version")
	}
	methods_len, _ := r.ReadByte()
	buf := make([]byte, methods_len)
	io.ReadFull(r, buf)
	resp := []byte{5, 0} //以上操作实现了接受客户端消息，所以服务器需要回应客户端消息。第一个参数表示版本号为5，即socks5协议，第二个参数表示认证方式为0，即无需密码访问。
	_, err := conn.Write(resp)
	return err
}

func Socks5ReadAddr(r *bufio.Reader) (string, error) {
	version, _ := r.ReadByte()
	if version != SOCKS5 {
		return "", errors.New("not SOCKS5 version")
	}
	cmd, _ := r.ReadByte()
	if cmd != CmdConnect {
		return "", errors.New("not SOCKS5 connect cmd")
	}
	r.ReadByte() //跳过RSV字段，即RSV保留字端，值长度为1个字节。
	addrtype, _ := r.ReadByte()
	var addr []byte
	switch addrtype {
	case AtypDomainName:
		log.Println("domain yes")
		addrlen, _ := r.ReadByte()   //读取一个字节以得到域名的长度。因为服务器地址类型的长度就是“1”，所以它是IP还是域名我们都能获取到完整的内容。如果能走到这一行代码说明一定是域名，如果没有上面的一行过滤代码我们就还需要考虑IPV4和IPV6的两种情况啦！
		addr = make([]byte, addrlen) //定义一个和域名长度一样大小的容器。
		io.ReadFull(r, addr)         //将域名的内容读取出来。
		break
	case AtypIPv4:
		return "", errors.New("SOCKS5 not support ipv4")
		break
	case AtypIPv6:
		return "", errors.New("SOCKS5 not support ipv6")
		break
	default:
		return "", errors.New("SOCKS5 not support type")
	}
	var port int16                          //因为端口是有2个字节来表示的，所以我们用int16来定义它的取值范围就OK。
	binary.Read(r, binary.BigEndian, &port) //读取2个字节，并将读取到的内容赋值给port变量。
	return fmt.Sprintf("%s:%d", addr, port), nil
}

func Socks5ConnectionConfirm(conn net.Conn) error {
	resp := []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00} //详情请参考：http://www.cnblogs.com/yinzhengjie/p/7357860.html
	_, err := conn.Write(resp)
	return err
}
