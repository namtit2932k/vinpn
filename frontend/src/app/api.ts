// The only module that imports generated Wails bindings. Tests mock it.
export * as Service from "../../bindings/github.com/sickyturtlez/vinpn/internal/app/service";
export type {
  AppError,
  AppInfo,
  AutotuneProgress,
  LogEvent,
  ScanProgress,
  ServerRow,
  Snapshot,
  StatsEvent,
  UpdateInfo,
  LANInfo,
  RulesView,
  RulesCompiled,
  ListsProgress,
  ProxyStatus,
  DeviceInfo,
  FakeSNIView,
  DNSServerStatus,
  FakeSNIStatus,
  SetupCountdown,
  LookupResult,
  StampCard,
  AdvRow,
  AdvScanProgress,
  AdvScanStart,
  ServerFilter,
  CFView,
  CFProgress,
  ImportPreview,
} from "../../bindings/github.com/sickyturtlez/vinpn/internal/app/models";
export type { Answer, Source as LookupSource, Verdict } from "../../bindings/github.com/sickyturtlez/vinpn/internal/lookup/models";
export type { Fields as StampFields } from "../../bindings/github.com/sickyturtlez/vinpn/internal/stamps/models";
export type { Result as AdvResult } from "../../bindings/github.com/sickyturtlez/vinpn/internal/scanner/advanced/models";
export type { Result as CfResult } from "../../bindings/github.com/sickyturtlez/vinpn/internal/cfscan/models";
export type { Choices as ImportChoices, Preview as BackupPreview, Warning as BackupWarning } from "../../bindings/github.com/sickyturtlez/vinpn/internal/backup/models";
export type { ServeStats } from "../../bindings/github.com/sickyturtlez/vinpn/internal/engine/models";
export type { Cert } from "../../bindings/github.com/sickyturtlez/vinpn/internal/certstore/models";
export type { Rule, Decision, LineError, Source as RuleSource } from "../../bindings/github.com/sickyturtlez/vinpn/internal/rules/models";
export type { List, CatalogItem } from "../../bindings/github.com/sickyturtlez/vinpn/internal/rules/lists/models";
export type { Stats as ProxyStats, ConnEvent } from "../../bindings/github.com/sickyturtlez/vinpn/internal/proxy/models";
export type { Settings, UpstreamProxy } from "../../bindings/github.com/sickyturtlez/vinpn/internal/store/models";
export type { QueryEvent } from "../../bindings/github.com/sickyturtlez/vinpn/internal/engine/models";
export type { Result as ProbeResult } from "../../bindings/github.com/sickyturtlez/vinpn/internal/probe/models";
export type { Adapter } from "../../bindings/github.com/sickyturtlez/vinpn/internal/sysdns/models";

export type StatusName = "disconnected" | "connecting" | "protected" | "degraded" | "disconnecting" | "error";
