//路径：src/utils/colors.ts
//默认可选颜色
export const defaultColors = ['#3491FA','#16C9FF','#34C388', '#6ECBEE','#F76965','#F76965','#FF9916', '#D08CF1', '#907AEF','#722ED1']
/**
 * 根据topicList数量获取下一个主题颜色
 * @param count topic列表个数
 * @returns color string
 */
export function getNextTopicColor(count:number) {
  const idx = count % defaultColors.length
  return defaultColors[idx]
}
// 生成随机颜色
export const getRandomColor = (): string => {
  const letters = '0123456789ABCDEF'
  let color = '#'
  for (let i = 0; i < 6; i += 1) {
    color += letters[Math.floor(Math.random() * 16)]
  }
  return color
}

