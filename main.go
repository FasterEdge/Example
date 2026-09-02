// ─────────────────────────────────────────────────────────────
// FasterEdge 开源项目
// Github: https://github.com/FasterEdge
// Gitee:  https://gitee.com/FasterEdge
// ─────────────────────────────────────────────────────────────
// Command example 是 FasterEdge 主项目生态的完整入门示例（demo）。
//
// 它演示如何通过公共 API 使用 github.com/FasterEdge/FasterEdge：
//  1. 初始化一个标准 Atom（InitStandardAtom）
//  2. 挂载全部内置能力（PreRun）
//  3. 调用各能力命令：时间同步 / 网络拓扑 / 一键令牌 / 命令执行
//  4. 演示"鉴权命令"链路（OneKey 签发令牌 → AuthenticatedCommand 校验）
//  5. 演示 Atom 完整生命周期（RunAtom 运行 → 优雅关闭）
//
// 用法：
//
//	go run .              # 运行 demo
//	go build .            # 编译为本地可执行文件
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/FasterEdge/FasterEdge"
	"github.com/FasterEdge/FasterEdge/ability"
	"github.com/FasterEdge/FasterEdge/data"
	"github.com/FasterEdge/FasterEdge/types"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("  FasterEdge Example Demo")
	fmt.Println("========================================")
	fmt.Printf("  运行时: %s (%s/%s)\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Println()

	// ---- 1. 初始化标准 Atom ----
	// InitStandardAtom 注册一整套常用 Data + Ability：
	// BaseData / NetMapData / KeyringData / ConfigData / 各数据库 Data，
	// RoleAbility / TimeAbility / NetMapAbility / OneKeyAbility /
	// CmdAbility / ShAbility / BashAbility / ConfigFileAbility，
	// 并自动把 OneKeyAbility 配置为命令鉴权器。
	atom := FasterEdge.InitStandardAtom()
	if err := atom.SetName("example-demo-atom"); err != nil {
		exit("设置 Atom 名称", err)
	}
	fmt.Printf("✅ Atom 已创建: %q (状态=%v, 能力=%d)\n",
		atom.GetName(), atom.State(), len(atom.AllAbilities()))

	// ---- 2. 挂载（PreRun）----
	// 只有处于 AtomCreated 状态的 Atom 可以调用 PreRun 进入 Mounted 状态。
	if err := FasterEdge.PreRunAtom(atom); err != nil {
		exit("挂载 Atom", err)
	}
	fmt.Printf("✅ Atom 已挂载 (状态=%v)\n", atom.State())

	// ---- 3. 能力命令演示 ----
	demoTimeAbility(atom)
	demoNetMapAbility(atom)
	demoOneKeyAbility(atom)
	demoCmdAbility(atom)
	demoAuthenticatedCommand(atom)

	// ---- 4. 生命周期：运行 → 优雅关闭 ----
	demoLifecycle(atom)

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("  Demo 完成，全部能力调用正常")
	fmt.Println("========================================")
}

// demoTimeAbility 演示时间能力：get_time 返回当前时间快照。
func demoTimeAbility(atom *types.Atom) {
	section("时间能力 TimeAbility")
	out := atom.CommandContext(context.Background(), "TimeAbility", ability.TimeCommandGetTime, nil)
	if out.Err != nil {
		fmt.Printf("  ⚠️  get_time 失败: %v\n", out.Err)
		return
	}
	ts, ok := out.Value.(ability.TimeSnapshot)
	if !ok {
		fmt.Printf("  ⚠️  get_time 返回值类型异常: %#v\n", out.Value)
		return
	}
	fmt.Printf("  ✅ 当前时间: %s (来源: %s)\n", ts.Time.Format(time.RFC3339), ts.Source)
}

// demoNetMapAbility 演示网络拓扑能力：注册对等节点并查询列表。
func demoNetMapAbility(atom *types.Atom) {
	section("网络拓扑能力 NetMapAbility")

	// 注册一个对等节点
	regOut := atom.CommandContext(context.Background(), "NetMapAbility",
		ability.NetMapCommandRegisterPeer,
		ability.NetMapRegisterPeerArgs{Name: "edge-node-01", Address: "192.168.1.100", Role: "edge"})
	if regOut.Err != nil {
		fmt.Printf("  ⚠️  register_peer 失败: %v\n", regOut.Err)
		return
	}
	fmt.Printf("  ✅ 已注册对等节点 edge-node-01\n")

	// 查询全部节点
	listOut := atom.CommandContext(context.Background(), "NetMapAbility",
		ability.NetMapCommandListPeers, nil)
	if listOut.Err != nil {
		fmt.Printf("  ⚠️  list_peers 失败: %v\n", listOut.Err)
		return
	}
	fmt.Printf("  ✅ 当前拓扑节点: %v\n", listOut.Value)
}

// demoOneKeyAbility 演示一键令牌能力：为对等节点签发短期访问令牌。
func demoOneKeyAbility(atom *types.Atom) {
	section("一键令牌能力 OneKeyAbility")

	out := atom.CommandContext(context.Background(), "OneKeyAbility",
		ability.OneKeyCommandIssueToken,
		ability.OneKeyIssueTokenArgs{Subject: "edge-node-01", TTL: time.Hour})
	if out.Err != nil {
		fmt.Printf("  ⚠️  issue_token 失败: %v\n", out.Err)
		return
	}
	tok, ok := out.Value.(ability.OneKeyToken)
	if !ok {
		fmt.Printf("  ⚠️  issue_token 返回值类型异常: %#v\n", out.Value)
		return
	}
	fmt.Printf("  ✅ 已为 edge-node-01 签发令牌: 有效期 %s → %s\n",
		tok.IssuedAt.Format(time.RFC3339), tok.ExpiresAt.Format(time.RFC3339))

	// 用 verify_token 在"远端"模拟校验
	verifyOut := atom.CommandContext(context.Background(), "OneKeyAbility",
		ability.OneKeyCommandVerifyToken,
		ability.OneKeyVerifyTokenArgs{
			Subject:   tok.Subject,
			IssuedAt:  tok.IssuedAt,
			ExpiresAt: tok.ExpiresAt,
			Signature: tok.Signature,
		})
	if verifyOut.Err != nil {
		fmt.Printf("  ⚠️  verify_token 失败: %v\n", verifyOut.Err)
		return
	}
	fmt.Printf("  ✅ 令牌校验通过: %v\n", verifyOut.Value)
}

// demoCmdAbility 演示命令执行能力：配置白名单并运行一条受控命令。
func demoCmdAbility(atom *types.Atom) {
	section("命令执行能力 CmdAbility")

	// 配置执行白名单：允许运行 go 与 echo（限制 go 最多 1 个参数）
	setOut := atom.CommandContext(context.Background(), "CmdAbility",
		ability.CmdCommandSetAllowlist,
		ability.CmdSetAllowlistArgs{Entries: []ability.CmdAllowlistEntry{
			{Name: "go", ArgsPrefix: []string{"version"}, MaxArgs: 1},
			{Name: "echo", ArgsPrefix: nil, MaxArgs: 3},
		}})
	if setOut.Err != nil {
		fmt.Printf("  ⚠️  set_allowlist 失败: %v\n", setOut.Err)
		return
	}
	fmt.Printf("  ✅ 执行白名单已配置: %v\n", setOut.Value)

	// 运行一条白名单内命令（Windows 下 echo 由 shell 提供，可能不可用，做容错）
	runOut := atom.CommandContext(context.Background(), "CmdAbility",
		ability.CmdCommandRun,
		ability.CmdRunArgs{Name: "echo", Args: []string{"Hello from FasterEdge demo"}, Timeout: 10 * time.Second})
	if runOut.Err != nil {
		fmt.Printf("  ⚠️  run 命令失败（跨平台差异，属正常）: %v\n", runOut.Err)
		return
	}
	if res, ok := runOut.Value.(ability.CmdResult); ok {
		fmt.Printf("  ✅ 命令执行完成: exit=%d stdout=%q stderr=%q\n",
			res.ExitCode, res.Stdout, res.Stderr)
	}
}

// demoAuthenticatedCommand 演示鉴权命令链路：
// OneKey 令牌 → Atom.AuthenticatedCommand 通过安全边界派发命令。
func demoAuthenticatedCommand(atom *types.Atom) {
	section("鉴权命令 AuthenticatedCommand")

	// 通过 OneKeyAbility 签发凭证（等价于远端持有令牌）
	issued := atom.CommandContext(context.Background(), "OneKeyAbility",
		ability.OneKeyCommandIssueToken,
		ability.OneKeyIssueTokenArgs{Subject: "remote-admin", TTL: 30 * time.Minute})
	if issued.Err != nil {
		fmt.Printf("  ⚠️ 签发令牌失败: %v\n", issued.Err)
		return
	}
	tok, _ := issued.Value.(ability.OneKeyToken)

	// 组装凭证并走鉴权入口查询 BaseData 信息
	credential := ability.OneKeyCredential{
		Subject:   tok.Subject,
		IssuedAt:  tok.IssuedAt,
		ExpiresAt: tok.ExpiresAt,
		Signature: tok.Signature,
	}
	out := atom.AuthenticatedCommand(credential, "BaseData", data.CommandInfo, nil)
	if out.Err != nil {
		fmt.Printf("  ⚠️  鉴权命令失败: %v\n", out.Err)
		return
	}
	fmt.Printf("  ✅ 鉴权命令执行成功（BaseData.%s）: %v\n", data.CommandInfo, out.Value)

	// 篡改签名应被拒绝 → 验证安全边界
	bad := credential
	bad.Signature += "tampered"
	if out := atom.AuthenticatedCommand(bad, "BaseData", data.CommandInfo, nil); out.Err != nil {
		fmt.Printf("  ✅ 篡改签名已被拒绝（安全边界生效）: %v\n", out.Err)
	} else {
		fmt.Printf("  ⚠️  篡改签名未被拒绝！")
	}
}

// demoLifecycle 演示 Atom 运行与优雅关闭。
func demoLifecycle(atom *types.Atom) {
	section("生命周期 RunAtom → CloseAtom")

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		// RunAll 会运行所有实现 Runner 的能力并阻塞，
		// 直到 ctx 取消后统一清理并回到 Stopped 状态。
		done <- FasterEdge.RunAtom(ctx, atom, FasterEdge.WithShutdownTimeout(3*time.Second))
	}()

	// 模拟运行一段时间
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("  ✅ Atom 运行中 (状态=%v)\n", atom.State())

	// 通知停止并等待优雅关闭。
	// context.Canceled 是取消 ctx 后的预期结果，属正常退出路径。
	cancel()
	if err := <-done; err != nil && !errors.Is(err, context.Canceled) {
		fmt.Printf("  ⚠️  RunAtom 返回错误: %v\n", err)
	} else {
		fmt.Printf("  ✅ Atom 已优雅停止 (状态=%v)\n", atom.State())
	}
}

func section(title string) {
	fmt.Println()
	fmt.Printf("── %s ──\n", title)
}

func exit(what string, err error) {
	fmt.Fprintf(os.Stderr, "❌ %s失败: %v\n", what, err)
	os.Exit(1)
}
