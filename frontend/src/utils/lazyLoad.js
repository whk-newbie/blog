/**
 * 图片懒加载工具
 * 使用Intersection Observer API实现图片懒加载
 */

/**
 * 创建懒加载图片指令
 */
export function setupLazyLoadDirective(app) {
  app.directive('lazy', {
    mounted(el, binding) {
      // 如果浏览器不支持Intersection Observer，直接加载图片
      if (!('IntersectionObserver' in window)) {
        el.src = binding.value
        return
      }

      // 设置占位符
      const placeholder = 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1 1"%3E%3C/svg%3E'
      el.src = placeholder
      el.style.backgroundColor = '#f5f5f5'
      el.style.minHeight = '200px'

      // 创建Intersection Observer
      const observer = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            if (entry.isIntersecting) {
              // 图片进入视口，开始加载
              const img = new Image()
              img.onload = () => {
                el.src = binding.value
                el.style.backgroundColor = 'transparent'
                el.style.minHeight = 'auto'
              }
              img.onerror = () => {
                // 加载失败，显示占位符
                el.src = placeholder
                el.alt = '图片加载失败'
              }
              img.src = binding.value
              observer.unobserve(el)
            }
          })
        },
        {
          rootMargin: '50px', // 提前50px开始加载
        }
      )

      observer.observe(el)
    },
    updated(el, binding) {
      // 如果图片URL改变，重新加载
      if (el.src !== binding.value) {
        el.src = binding.value
      }
    },
  })
}

/**
 * 懒加载图片组件
 */
export function createLazyImage(src, alt = '', className = '') {
  return {
    template: `
      <img
        v-lazy="src"
        :alt="alt"
        :class="className"
        loading="lazy"
      />
    `,
    props: {
      src: {
        type: String,
        required: true,
      },
      alt: {
        type: String,
        default: '',
      },
      className: {
        type: String,
        default: '',
      },
    },
  }
}

