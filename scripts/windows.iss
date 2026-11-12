; MyClip per-user Windows installer, built by CI with Inno Setup.
; The version comes from mygo.json, passed as /DAppVersion by the caller.

#define AppName "MyClip"
#ifndef AppVersion
#error Pass /DAppVersion=x.y.z (the mygo.json version)
#endif

[Setup]
AppId={{016570d7-f533-4b46-b703-d6fab5872225}
AppName={#AppName}
AppVersion={#AppVersion}
AppPublisher=L92-Labs
DefaultDirName={autopf}\{#AppName}
PrivilegesRequired=lowest
DisableProgramGroupPage=yes
OutputDir=..\build\windows-amd64
OutputBaseFilename=myclip-setup
Compression=lzma2
SolidCompression=yes
UninstallDisplayIcon={app}\myclip.exe
CloseApplications=yes

[Files]
Source: "..\build\windows-amd64\myclip.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\LICENSE"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\NOTICE.md"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\THIRD_PARTY_NOTICES.txt"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{autoprograms}\{#AppName}"; Filename: "{app}\myclip.exe"
Name: "{autodesktop}\{#AppName}"; Filename: "{app}\myclip.exe"; Tasks: desktopicon

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Run]
Filename: "{app}\myclip.exe"; Description: "{cm:LaunchProgram,{#AppName}}"; Flags: nowait postinstall skipifsilent
