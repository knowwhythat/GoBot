import { h, inject, computed } from "vue";
import {
  riH1,
  riH2,
  riAB,
  riBold,
  riLink,
  riLinkM,
  riItalic,
  riTable2,
  riEyeLine,
  riAddLine,
  riSortAsc,
  riMindMap,
  riKey2Line,
  riTBoxLine,
  riSaveLine,
  riPlayLine,
  riPlayCircleLine,
  riMoreLine,
  riStopLine,
  riSortDesc,
  riTimeLine,
  riFlagLine,
  riFileLine,
  riTeamLine,
  riLinksLine,
  riGroupLine,
  riGuideLine,
  riChat3Line,
  riEarthLine,
  riLock2Line,
  riSliceLine,
  riHome5Line,
  riShareLine,
  riBook3Line,
  riPauseLine,
  riFlowChart,
  riMore2Line,
  riMouseLine,
  riFocusLine,
  riFontSize2,
  riParagraph,
  riImageLine,
  riCloseLine,
  riCheckLine,
  riTimerLine,
  riMagicLine,
  riHtml5Line,
  riToggleFill,
  riToggleLine,
  riFolderLine,
  riAlertLine,
  riGithubFill,
  riEyeOffLine,
  riWindowLine,
  riPencilLine,
  riBrush2Line,
  riGlobalLine,
  riShieldLine,
  riCursorLine,
  riUploadLine,
  riFocus3Line,
  riTwitterLine,
  riDiscordLine,
  riLinkUnlinkM,
  riYoutubeLine,
  riSideBarLine,
  riSideBarFill,
  riWindow2Line,
  riRefreshLine,
  riRefreshFill,
  riFilter2Line,
  riRestartLine,
  riSearch2Line,
  riSendPlaneLine,
  riEditBoxLine,
  riHistoryLine,
  riRepeat2Line,
  riCommandLine,
  riArticleLine,
  riKeyboardLine,
  riFileEditLine,
  riCompass3Line,
  riFolderOpenLine,
  riComputerLine,
  riFileCopyLine,
  riCalendarLine,
  riFileTextLine,
  riSubtractLine,
  riBracketsLine,
  riPushpin2Line,
  riPushpin2Fill,
  riDownloadCloud2Line,
  riDownloadLine,
  riFileListLine,
  riDragDropLine,
  riClipboardLine,
  riCheckDoubleLine,
  riDoubleQuotesL,
  riLightbulbLine,
  riFolderZipLine,
  riFileShredLine,
  riHandHeartLine,
  riDatabase2Line,
  riSettings3Line,
  riArrowLeftLine,
  riEqualizerLine,
  riStrikethrough2,
  riFileUploadLine,
  riCodeSSlashLine,
  riHardDrive2Line,
  riDeleteBin7Line,
  riArrowLeftSLine,
  riFullscreenLine,
  riFlashlightLine,
  riTimerFlashLine,
  riBaseStationLine,
  riFileHistoryLine,
  riInformationLine,
  riArrowUpDownLine,
  riArrowGoBackLine,
  riInputCursorMove,
  riCloseCircleLine,
  riRecordCircleLine,
  riErrorWarningLine,
  riExternalLinkLine,
  riUploadCloud2Line,
  riFileDownloadLine,
  riShieldKeyholeLine,
  riArrowDropDownLine,
  riNotification3Line,
  riArrowLeftRightLine,
  riArrowGoForwardLine,
  riCheckboxCircleLine,
  riLightbulbFlashLine,
  riIncreaseDecreaseLine,
  riPriceTag3Line,
  riArrowLeftDownFill,
  riArrowLeftCircleLine,
  riCodepenLine,
  riReactjsFill,
  riCompassesLine,
  riReplyAllLine,
  riGitCommitLine,
  riGitRepositoryLine,
  riBugLine,
  riSpeedMiniLine,
  riSkipForwardMiniLine,
  riChromeLine,
  riPagesLine,
  riMenuAddLine,
  riSwitchLine,
  riScreenshotLine,
  riShapeLine,
  riArrowRightCircleLine,
  riText,
  riFridgeLine,
  riTranslate2,
  riCharacterRecognitionLine,
  riLockPasswordLine,
  riChatSettingsLine,
  riChatNewLine,
  riChatDeleteLine,
  riIndentDecrease,
  riIndentIncrease,
  riTerminalWindowLine,
  riLayoutBottomLine,
  riUninstallLine,
  riCrop2Line,
  riCheckboxLine,
  riFocus2Line,
  riNodeTree,
} from "v-remixicon/icons";

export const icons = {
  riH1,
  riH2,
  riAB,
  riBold,
  riLink,
  riLinkM,
  riItalic,
  riTable2,
  riEyeLine,
  riAddLine,
  riSortAsc,
  riMindMap,
  riKey2Line,
  riTBoxLine,
  riSaveLine,
  riPlayLine,
  riPlayCircleLine,
  riMoreLine,
  riStopLine,
  riSortDesc,
  riTimeLine,
  riFlagLine,
  riFileLine,
  riTeamLine,
  riLinksLine,
  riGroupLine,
  riGuideLine,
  riChat3Line,
  riEarthLine,
  riLock2Line,
  riSliceLine,
  riHome5Line,
  riShareLine,
  riBook3Line,
  riPauseLine,
  riFlowChart,
  riMore2Line,
  riMouseLine,
  riFocusLine,
  riFontSize2,
  riParagraph,
  riImageLine,
  riCloseLine,
  riCheckLine,
  riTimerLine,
  riMagicLine,
  riHtml5Line,
  riToggleFill,
  riToggleLine,
  riFolderLine,
  riAlertLine,
  riGithubFill,
  riEyeOffLine,
  riWindowLine,
  riPencilLine,
  riBrush2Line,
  riGlobalLine,
  riShieldLine,
  riCursorLine,
  riUploadLine,
  riFocus3Line,
  riTwitterLine,
  riDiscordLine,
  riLinkUnlinkM,
  riYoutubeLine,
  riSideBarLine,
  riSideBarFill,
  riWindow2Line,
  riRefreshLine,
  riRefreshFill,
  riFilter2Line,
  riRestartLine,
  riSearch2Line,
  riSendPlaneLine,
  riEditBoxLine,
  riHistoryLine,
  riRepeat2Line,
  riCommandLine,
  riArticleLine,
  riKeyboardLine,
  riFileEditLine,
  riCompass3Line,
  riFolderOpenLine,
  riComputerLine,
  riFileCopyLine,
  riCalendarLine,
  riFileTextLine,
  riSubtractLine,
  riBracketsLine,
  riPushpin2Line,
  riPushpin2Fill,
  riDownloadCloud2Line,
  riDownloadLine,
  riFileListLine,
  riDragDropLine,
  riClipboardLine,
  riCheckDoubleLine,
  riDoubleQuotesL,
  riLightbulbLine,
  riFolderZipLine,
  riFileShredLine,
  riHandHeartLine,
  riDatabase2Line,
  riSettings3Line,
  riArrowLeftLine,
  riEqualizerLine,
  riStrikethrough2,
  riFileUploadLine,
  riCodeSSlashLine,
  riHardDrive2Line,
  riDeleteBin7Line,
  riArrowLeftSLine,
  riFullscreenLine,
  riFlashlightLine,
  riTimerFlashLine,
  riBaseStationLine,
  riFileHistoryLine,
  riInformationLine,
  riArrowUpDownLine,
  riArrowGoBackLine,
  riInputCursorMove,
  riCloseCircleLine,
  riRecordCircleLine,
  riErrorWarningLine,
  riExternalLinkLine,
  riUploadCloud2Line,
  riFileDownloadLine,
  riShieldKeyholeLine,
  riArrowDropDownLine,
  riNotification3Line,
  riArrowLeftRightLine,
  riArrowGoForwardLine,
  riCheckboxCircleLine,
  riLightbulbFlashLine,
  riIncreaseDecreaseLine,
  riPriceTag3Line,
  riArrowLeftDownFill,
  riArrowLeftCircleLine,
  riCodepenLine,
  riReactjsFill,
  riCompassesLine,
  riReplyAllLine,
  riGitCommitLine,
  riGitRepositoryLine,
  riBugLine,
  riSpeedMiniLine,
  riSkipForwardMiniLine,
  riChromeLine,
  riPagesLine,
  riMenuAddLine,
  riSwitchLine,
  riScreenshotLine,
  riShapeLine,
  riArrowRightCircleLine,
  riText,
  riFridgeLine,
  riTranslate2,
  riCharacterRecognitionLine,
  riLockPasswordLine,
  riChatSettingsLine,
  riChatNewLine,
  riChatDeleteLine,
  riIndentDecrease,
  riIndentIncrease,
  riTerminalWindowLine,
  riLayoutBottomLine,
  riUninstallLine,
  riCrop2Line,
  riCheckboxLine,
  riFocus2Line,
  riNodeTree,
  mdiEqual: "M19,10H5V8H19V10M19,16H5V14H19V16Z",
  mdiPackageVariantClosed:
    "M21,16.5C21,16.88 20.79,17.21 20.47,17.38L12.57,21.82C12.41,21.94 12.21,22 12,22C11.79,22 11.59,21.94 11.43,21.82L3.53,17.38C3.21,17.21 3,16.88 3,16.5V7.5C3,7.12 3.21,6.79 3.53,6.62L11.43,2.18C11.59,2.06 11.79,2 12,2C12.21,2 12.41,2.06 12.57,2.18L20.47,6.62C20.79,6.79 21,7.12 21,7.5V16.5M12,4.15L10.11,5.22L16,8.61L17.96,7.5L12,4.15M6.04,7.5L12,10.85L13.96,9.75L8.08,6.35L6.04,7.5M5,15.91L11,19.29V12.58L5,9.21V15.91M19,15.91V9.21L13,12.58V19.29L19,15.91Z",
  mdiVariable:
    "M20.41,3C21.8,5.71 22.35,8.84 22,12C21.8,15.16 20.7,18.29 18.83,21L17.3,20C18.91,17.57 19.85,14.8 20,12C20.34,9.2 19.89,6.43 18.7,4L20.41,3M5.17,3L6.7,4C5.09,6.43 4.15,9.2 4,12C3.66,14.8 4.12,17.57 5.3,20L3.61,21C2.21,18.29 1.65,15.17 2,12C2.2,8.84 3.3,5.71 5.17,3M12.08,10.68L14.4,7.45H16.93L13.15,12.45L15.35,17.37H13.09L11.71,14L9.28,17.33H6.76L10.66,12.21L8.53,7.45H10.8L12.08,10.68Z",
  mdiRegex:
    "M16,16.92C15.67,16.97 15.34,17 15,17C14.66,17 14.33,16.97 14,16.92V13.41L11.5,15.89C11,15.5 10.5,15 10.11,14.5L12.59,12H9.08C9.03,11.67 9,11.34 9,11C9,10.66 9.03,10.33 9.08,10H12.59L10.11,7.5C10.3,7.25 10.5,7 10.76,6.76V6.76C11,6.5 11.25,6.3 11.5,6.11L14,8.59V5.08C14.33,5.03 14.66,5 15,5C15.34,5 15.67,5.03 16,5.08V8.59L18.5,6.11C19,6.5 19.5,7 19.89,7.5L17.41,10H20.92C20.97,10.33 21,10.66 21,11C21,11.34 20.97,11.67 20.92,12H17.41L19.89,14.5C19.7,14.75 19.5,15 19.24,15.24V15.24C19,15.5 18.75,15.7 18.5,15.89L16,13.41V16.92H16V16.92M5,19A2,2 0 0,1 7,17A2,2 0 0,1 9,19A2,2 0 0,1 7,21A2,2 0 0,1 5,19H5Z",
  mdiCookieOutline:
    "M20.87 10.5C20.6 10 20 10 20 10H18V9C18 8 17 8 17 8H15V7C15 6 14 6 14 6H13V4C13 3 12 3 12 3C7.03 3 3 7.03 3 12C3 16.97 7.03 21 12 21C16.97 21 21 16.97 21 12C21 11.5 20.96 11 20.87 10.5M11.32 18.96C12 18.82 12.5 18.22 12.5 17.5C12.5 16.67 11.83 16 11 16S9.5 16.67 9.5 17.5C9.5 18 9.76 18.47 10.16 18.74C7.54 18.04 5.5 15.81 5.09 13.12C5 12.61 5 12.11 5 11.62C5.07 12.39 5.71 13 6.5 13C7.33 13 8 12.33 8 11.5S7.33 10 6.5 10C5.82 10 5.25 10.46 5.07 11.08C5.47 8 7.91 5.5 11 5.07V6.5C11 7.33 11.67 8 12.5 8H13V8.5C13 9.33 13.67 10 14.5 10H16V10.5C16 11.33 16.67 12 17.5 12H19C19 16.08 15.5 19.36 11.32 18.96M9.5 9C8.67 9 8 8.33 8 7.5S8.67 6 9.5 6 11 6.67 11 7.5 10.33 9 9.5 9M13 12.5C13 13.33 12.33 14 11.5 14S10 13.33 10 12.5 10.67 11 11.5 11 13 11.67 13 12.5M18 14.5C18 15.33 17.33 16 16.5 16S15 15.33 15 14.5 15.67 13 16.5 13 18 13.67 18 14.5Z",
  mdiDrag:
    "M7,19V17H9V19H7M11,19V17H13V19H11M15,19V17H17V19H15M7,15V13H9V15H7M11,15V13H13V15H11M15,15V13H17V15H15M7,11V9H9V11H7M11,11V9H13V11H11M15,11V9H17V11H15M7,7V5H9V7H7M11,7V5H13V7H11M15,7V5H17V7H15Z",
  webhookIcon:
    "M11.1339746,6.5 C11.410117,6.02170738 12.0217074,5.85783222 12.5,6.1339746 C12.6618283,6.22740623 12.7876628,6.35923925 12.8726385,6.51131768 L12.8728316,6.51197847 L15.5280985,11.258439 C17.7104719,10.5794609 20.1487626,11.4726335 21.3395353,13.5351122 C22.7202472,15.9265753 21.9008714,18.9845274 19.5094083,20.3652392 C18.0316717,21.2184108 16.2512131,21.2502691 14.7625925,20.5024023 L14.5614635,20.3955914 L15.5391998,18.650876 C16.4572736,19.1653634 17.5817121,19.1687941 18.5094083,18.6331884 C19.9442862,17.8047613 20.4359116,15.9699901 19.6074845,14.5351122 C18.8086441,13.15148 17.0740576,12.6449281 15.6646114,13.35331 L15.5094083,13.437036 L15.4923149,13.408 L14.6887376,13.8579251 L11.1271684,7.48802153 C10.9611644,7.19043272 10.9514715,6.81610467 11.1339746,6.5 Z M3.69009425,12.2357175 L5.00912623,13.7390986 C4.15631502,14.4873355 3.78941522,15.6581523 4.08821537,16.7732896 C4.5170408,18.373688 6.16205167,19.3234354 7.76244998,18.89461 C9.00692272,18.5611545 9.88400974,17.477653 9.97851857,16.2240731 L9.98665134,16.0438959 L10.0023149,16.043 L10.0029449,15 L17,15 C17.1696897,14.9996082 17.342108,15.0428156 17.5,15.1339746 C17.9782926,15.410117 18.1421678,16.0217074 17.8660254,16.5 C17.6809741,16.8205182 17.3452819,16.9998396 17.000012,17.0001665 L17,17 L11.8855308,17.0005501 C11.512912,18.8188035 10.1440497,20.3270146 8.28008807,20.8264616 C5.61275756,21.5411707 2.87107277,19.9582583 2.15636371,17.2909277 C1.67853507,15.507647 2.22314594,13.6374116 3.52361777,12.3885391 L3.69009425,12.2357175 Z M12,2 C14.6887547,2 16.8818181,4.12230671 16.9953805,6.78311038 L17,7 L15,7 C15,5.34314575 13.6568542,4 12,4 C10.3431458,4 9,5.34314575 9,7 C9,8.03663591 9.52964394,8.97937227 10.3797433,9.52549875 L10.4663149,9.577 L11.4242264,10.1109289 L7.87365874,16.4865392 L7.8723149,16.486 L7.8660254,16.5 C7.60960748,16.9441289 7.06395121,17.1171529 6.60436072,16.9185096 L6.5,16.8660254 C6.02170738,16.589883 5.85783222,15.9782926 6.1339746,15.5 L6.12634126,15.5134608 L8.7500113,10.80028 C7.65887804,9.86710067 7,8.49115157 7,7 C7,4.23857625 9.23857625,2 12,2 Z",
  mdiGoogleSheet:
    "M19,11V9H11V5H9V9H5V11H9V19H11V11H19M19,3C19.5,3 20,3.2 20.39,3.61C20.8,4 21,4.5 21,5V19C21,19.5 20.8,20 20.39,20.39C20,20.8 19.5,21 19,21H5C4.5,21 4,20.8 3.61,20.39C3.2,20 3,19.5 3,19V5C3,4.5 3.2,4 3.61,3.61C4,3.2 4.5,3 5,3H19Z",
  mdiCursorDefaultClickOutline:
    "M11.5,11L17.88,16.37L17,16.55L16.36,16.67C15.73,16.8 15.37,17.5 15.65,18.07L15.92,18.65L17.28,21.59L15.86,22.25L14.5,19.32L14.24,18.74C13.97,18.15 13.22,17.97 12.72,18.38L12.21,18.78L11.5,19.35V11M10.76,8.69A0.76,0.76 0 0,0 10,9.45V20.9C10,21.32 10.34,21.66 10.76,21.66C10.95,21.66 11.11,21.6 11.24,21.5L13.15,19.95L14.81,23.57C14.94,23.84 15.21,24 15.5,24C15.61,24 15.72,24 15.83,23.92L18.59,22.64C18.97,22.46 19.15,22 18.95,21.63L17.28,18L19.69,17.55C19.85,17.5 20,17.43 20.12,17.29C20.39,16.97 20.35,16.5 20,16.21L11.26,8.86L11.25,8.87C11.12,8.76 10.95,8.69 10.76,8.69M15,10V8H20V10H15M13.83,4.76L16.66,1.93L18.07,3.34L15.24,6.17L13.83,4.76M10,0H12V5H10V0M3.93,14.66L6.76,11.83L8.17,13.24L5.34,16.07L3.93,14.66M3.93,3.34L5.34,1.93L8.17,4.76L6.76,6.17L3.93,3.34M7,10H2V8H7V10",
  restoreIcon:
    "M716.018 332.909H247.512c-50.729 0-92 41.271-92 92v321.54c0 50.729 41.271 92 92 92h468.506c50.729 0 92-41.271 92-92v-321.54c0-50.729-41.272-92-92-92z m36 413.54c0 19.851-16.149 36-36 36H247.512c-19.851 0-36-16.149-36-36v-321.54c0-19.851 16.149-36 36-36h468.506c19.851 0 36 16.149 36 36v321.54z M748.373 231.47H291.868c-15.464 0-28 12.536-28 28s12.536 28 28 28h456.505c61.757 0 112 50.243 112 112v309.54c0 15.464 12.536 28 28 28s28-12.536 28-28V399.47c0-92.636-75.364-168-168-168z",
  maxIcon:
    "M802.16 191.78H222c-16.57 0-30 12.22-30 27.28v585.39c0 15.06 13.43 27.27 30 27.27h585.16a25 25 0 0 0 25-25V219.06c0-15.06-13.43-27.28-30-27.28zM256.05 773.61V250.13h512.1v523.48z",
  minIcon:
    "M930.355 551.424H94.218c-23.987 0-43.975-17.408-43.975-39.424 0-21.504 19.42-39.424 43.975-39.424h835.564c23.987 0 43.975 17.408 43.975 39.424 0.006 22.016-19.415 39.424-43.402 39.424z",
  closeIcon:
    "M822.00345 776.822434l0.022513-0.022513L246.50423 201.317075c-5.78782-5.791913-13.785981-9.374508-22.621207-9.374508-17.662265 0-31.980365 14.3181-31.980365 31.980365 0 8.834202 3.582595 16.832364 9.373485 22.620184L776.11226 821.339324c5.838985 6.277984 14.166651 10.209526 23.416316 10.209526 17.662265 0 31.980365-14.3181 31.980365-31.980365C831.508941 790.667767 827.871087 782.620487 822.00345 776.822434z M776.783549 201.448058l-0.022513-0.022513L201.278189 776.947278c-5.791913 5.78782-9.374508 13.785981-9.374508 22.621207 0 17.662265 14.3181 31.980365 31.980365 31.980365 8.834202 0 16.832364-3.582595 22.620184-9.373485l574.797231-574.836117c6.277984-5.838985 10.209526-14.166651 10.209526-23.416316 0-17.662265-14.3181-31.980365-31.980365-31.980365C790.628882 191.942567 782.580578 195.58042 776.783549 201.448058z",
};

const component = {
  name: "v-remixicon",
  props: {
    name: String,
    title: String,
    viewBox: {
      type: String,
      default: "0 0 24 24",
    },
    size: {
      type: [String, Number],
      default: 24,
    },
    fill: {
      type: String,
      default: "currentColor",
    },
    rotate: {
      type: [Number, String],
      default: 0,
    },
    path: {
      type: String,
      default: "",
    },
  },
  setup(props) {
    const injectIcons = inject("remixicons");
    const icon = computed(() => {
      if (props.path) return props.path;

      const iconStr = injectIcons[props.name];

      if (typeof iconStr === "undefined") {
        console.error(
          `[v-remixicon] ${props.name} name of the icon is incorrect`
        );
        return null;
      }

      return iconStr;
    });

    return () =>
      h(
        "svg",
        {
          viewBox: props.viewBox,
          fill: props.fill,
          height: props.size,
          width: props.size,
          class: "v-remixicon",
          xmlns: "http://www.w3.org/2000/svg",
          style: {
            transform: props.rotate ? `rotate(${props.rotate}deg)` : null,
            display: "inline-block",
          },
        },
        [
          props.title ? h("title", {}, props.title) : null,
          h("g", {}, [
            h("path", { fill: "none", d: "M0 0h24v24H0z" }),
            h("path", { "fill-rule": "nonzero", d: icon.value }),
          ]),
        ]
      );
  },
};

export default {
  install(app) {
    app.provide("remixicons", icons);
    app.component("VRemixicon", component);
  },
};
