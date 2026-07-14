Add-Type -AssemblyName System.Drawing

$w = 740; $h = 420
$bgColor = '#0d0d1a'
$promptColor = '#00ff88'
$outputColor = '#888888'
$textColor = '#e0e0e0'
$titleBg = '#1a1a2e'

$font = New-Object System.Drawing.Font('Consolas', 12)
$promptBrush = New-Object System.Drawing.SolidBrush('#00ff88')
$outputBrush = New-Object System.Drawing.SolidBrush('#888888')
$textBrush = New-Object System.Drawing.SolidBrush('#e0e0e0')
$bgBrush = New-Object System.Drawing.SolidBrush('#0d0d1a')
$titleBrush = New-Object System.Drawing.SolidBrush('#666666')

$framesDir = "$env:TEMP\tc-frames"
New-Item -ItemType Directory -Path $framesDir -Force | Out-Null

$steps = @(
    @("PS> taskcapsule start feature-checkout --no-services",
      "Capsule started: feature-checkout",
      "Branch:    task/feature-checkout",
      "Status:    running"),
    @("PS> taskcapsule note feature-checkout " + '"Continue retry test next"',
      "Note saved for feature-checkout."),
    @("PS> taskcapsule pause feature-checkout",
      "Capsule paused: feature-checkout",
      "Resources released."),
    @("PS> taskcapsule start urgent-hotfix --no-services",
      "Capsule started: urgent-hotfix"),
    @("PS> taskcapsule pause urgent-hotfix",
      "Capsule paused: urgent-hotfix",
      "Resources released."),
    @("PS> taskcapsule resume feature-checkout",
      "Capsule resumed: feature-checkout",
      "Last note:  Continue retry test next"),
    @("PS> taskcapsule where feature-checkout",
      "You were working on: feature-checkout",
      "Last note:  Continue retry test next")
)

$allOutput = @()
$frameIndex = 0

foreach ($step in $steps) {
    $allOutput += $step
    
    # Generate 3 frames per step for animation timing
    for ($dup = 0; $dup -lt 3; $dup++) {
        $bmp = New-Object System.Drawing.Bitmap($w, $h)
        $g = [System.Drawing.Graphics]::FromImage($bmp)
        $g.Clear('#0d0d1a')
        
        # Title bar
        $rect = New-Object System.Drawing.Rectangle(0, 0, $w, 30)
        $titleBgBrush = New-Object System.Drawing.SolidBrush('#1a1a25')
        $g.FillRectangle($titleBgBrush, $rect)
        $g.DrawString("PowerShell - taskcapsule demo", $font, $titleBrush, 40, 6)
        
        # Title bar dots
        $g.FillEllipse([System.Drawing.Brushes.Red], 12, 10, 10, 10)
        $g.FillEllipse([System.Drawing.Brushes.Orange], 26, 10, 10, 10)
        $g.FillEllipse([System.Drawing.Brushes.LimeGreen], 40, 10, 10, 10)
        
        $y = 45
        $accumulated = @()
        foreach ($a in $allOutput) {
            $isPrompt = $a -and $a.StartsWith("PS>")
            if ($isPrompt) {
                $g.DrawString($a, $font, $promptBrush, 15, $y)
            } else {
                $g.DrawString($a, $font, $outputBrush, 15, $y)
            }
            $y += 22
        }
        
        $path = "$framesDir\frame_$("{0:D4}" -f $frameIndex).png"
        $bmp.Save($path)
        $g.Dispose()
        $bmp.Dispose()
        $frameIndex++
    }
}

# Generate final resting frame (more dupes for longer display)
$bmp = New-Object System.Drawing.Bitmap($w, $h)
$g = [System.Drawing.Graphics]::FromImage($bmp)
$g.Clear('#0d0d1a')
$titleBgBrush = New-Object System.Drawing.SolidBrush('#1a1a25')
$g.FillRectangle($titleBgBrush, 0, 0, $w, 30)
$g.DrawString("PowerShell - taskcapsule demo", $font, $titleBrush, 40, 6)
$g.FillEllipse([System.Drawing.Brushes.Red], 12, 10, 10, 10)
$g.FillEllipse([System.Drawing.Brushes.Orange], 26, 10, 10, 10)
$g.FillEllipse([System.Drawing.Brushes.LimeGreen], 40, 10, 10, 10)
$y = 45
foreach ($a in $allOutput) {
    if ($a.StartsWith("PS>")) { $g.DrawString($a, $font, $promptBrush, 15, $y) }
    else { $g.DrawString($a, $font, $outputBrush, 15, $y) }
    $y += 22
}
for ($dup = 0; $dup -lt 10; $dup++) {
    $path = "$framesDir\frame_$("{0:D4}" -f $frameIndex).png"
    $bmp.Save($path); $frameIndex++
}
$g.Dispose(); $bmp.Dispose()

# Use ffmpeg to create GIF
$outputPath = "$pwd\docs\assets\taskcapsule-demo.gif"
& ffmpeg -y -framerate 8 -i "$framesDir\frame_%04d.png" -filter_complex "[0:v] fps=8,split[s0][s1];[s0]palettegen=max_colors=128[p];[s1][p]paletteuse=dither=bayer" $outputPath 2>&1 | Out-Null

Write-Output "GIF created: $outputPath"
Write-Output "Total frames: $frameIndex"
Remove-Item -Recurse -Force $framesDir -ErrorAction SilentlyContinue
