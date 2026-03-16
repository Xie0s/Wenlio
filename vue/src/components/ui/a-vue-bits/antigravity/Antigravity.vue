<!--
  Antigravity 粒子动画组件
  
  功能说明：
  - 基于 Three.js 的 WebGL 粒子效果
  - 粒子跟随鼠标形成环形磁场效果
  - 支持自动动画模式
  - 可配置粒子形状、颜色、大小等参数
  
  参考来源：vue-bits (https://vue-bits.dev/animations/antigravity)
-->
<script setup lang="ts">
import { onMounted, onUnmounted, useTemplateRef, watch, ref, computed } from 'vue';
import {
  BoxGeometry,
  Clock,
  CylinderGeometry,
  InstancedMesh,
  MeshBasicMaterial,
  Object3D,
  PerspectiveCamera,
  Scene,
  SphereGeometry,
  TetrahedronGeometry,
  WebGLRenderer,
} from 'three';
import type { BufferGeometry, Material } from 'three';

export type ParticleShape = 'capsule' | 'sphere' | 'box' | 'tetrahedron';

interface AntigravityProps {
  count?: number;
  magnetRadius?: number;
  ringRadius?: number;
  waveSpeed?: number;
  waveAmplitude?: number;
  particleSize?: number;
  lerpSpeed?: number;
  color?: string;
  autoAnimate?: boolean;
  particleVariance?: number;
  rotationSpeed?: number;
  depthFactor?: number;
  pulseSpeed?: number;
  particleShape?: ParticleShape;
  fieldStrength?: number;
}

interface Particle {
  t: number;
  factor: number;
  speed: number;
  xFactor: number;
  yFactor: number;
  zFactor: number;
  mx: number;
  my: number;
  mz: number;
  nx: number; // Normalized X (-0.5 to 0.5)
  ny: number; // Normalized Y (-0.5 to 0.5)
  cx: number;
  cy: number;
  cz: number;
  vx: number;
  vy: number;
  vz: number;
  randomRadiusOffset: number;
}

const props = withDefaults(defineProps<AntigravityProps>(), {
  count: 300,
  magnetRadius: 10,
  ringRadius: 10,
  waveSpeed: 0.4,
  waveAmplitude: 1,
  particleSize: 2,
  lerpSpeed: 0.1,
  color: '#27FF64',
  autoAnimate: false,
  particleVariance: 1,
  rotationSpeed: 0,
  depthFactor: 1,
  pulseSpeed: 3,
  particleShape: 'capsule',
  fieldStrength: 10
});

const containerRef = useTemplateRef<HTMLDivElement>('containerRef');
const webglAvailable = ref(true);

let renderer: WebGLRenderer | null = null;
let scene: Scene | null = null;
let camera: PerspectiveCamera | null = null;
let mesh: InstancedMesh | null = null;
let animationFrameId: number = 0;
let particles: Particle[] = [];
let dummy: Object3D;
let lastMousePos = { x: 0, y: 0 };
let lastMouseMoveTime = 0;
let virtualMouse = { x: 0, y: 0 };
let pointer = { x: 0, y: 0 };
let clock: Clock;
let currentViewport = { width: 0, height: 0 };
let isVisible = true;

const fallbackStyle = computed(() => {
  if (webglAvailable.value) return {};
  // Simple fallback: a subtle gradient using the primary color
  return {
    background: `radial-gradient(circle at center, ${props.color}10 0%, transparent 70%)`
  };
});

function createGeometry(shape: ParticleShape): BufferGeometry {
  switch (shape) {
    case 'sphere':
      return new SphereGeometry(0.2, 16, 16);
    case 'box':
      return new BoxGeometry(0.3, 0.3, 0.3);
    case 'tetrahedron':
      return new TetrahedronGeometry(0.3);
    case 'capsule':
    default:
      return new CylinderGeometry(0.1, 0.1, 0.6, 12);
  }
}

function initParticles() {
  particles = [];
  for (let i = 0; i < props.count; i++) {
    const t = Math.random() * 100;
    const factor = 20 + Math.random() * 100;
    const speed = 0.01 + Math.random() / 200;
    const xFactor = -50 + Math.random() * 100;
    const yFactor = -50 + Math.random() * 100;
    const zFactor = -50 + Math.random() * 100;

    // Store normalized position (-0.5 to 0.5)
    const nx = Math.random() - 0.5;
    const ny = Math.random() - 0.5;
    
    const z = (Math.random() - 0.5) * 20;

    const randomRadiusOffset = (Math.random() - 0.5) * 2;

    // Calculate initial absolute position based on current viewport
    const x = nx * currentViewport.width;
    const y = ny * currentViewport.height;

    particles.push({
      t,
      factor,
      speed,
      xFactor,
      yFactor,
      zFactor,
      mx: x,
      my: y,
      mz: z,
      nx,
      ny,
      cx: x,
      cy: y,
      cz: z,
      vx: 0,
      vy: 0,
      vz: 0,
      randomRadiusOffset
    });
  }
}

function getViewportAtDepth(camera: PerspectiveCamera, depth: number) {
  const fovInRadians = (camera.fov * Math.PI) / 180;
  const height = 2 * Math.tan(fovInRadians / 2) * depth;
  const width = height * camera.aspect;
  return { width, height };
}

function isWebGLAvailable() {
  try {
    const canvas = document.createElement('canvas');
    // Basic WebGL availability check (allow software rendering)
    return !!(window.WebGLRenderingContext && 
      (canvas.getContext('webgl') || canvas.getContext('experimental-webgl'))
    );
  } catch (e) {
    return false;
  }
}

function setupScene() {
  const container = containerRef.value;
  if (!container) return;

  if (!isWebGLAvailable()) {
    webglAvailable.value = false;
    return;
  }

  try {
    const { clientWidth, clientHeight } = container;

    renderer = new WebGLRenderer({ 
      antialias: true,
      alpha: true,
      powerPreference: 'low-power'
    });
    renderer.setSize(clientWidth, clientHeight);
    renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
    container.appendChild(renderer.domElement);

    scene = new Scene();

    camera = new PerspectiveCamera(35, clientWidth / clientHeight, 0.1, 1000);
    camera.position.z = 50;

    // Initialize viewport size
    currentViewport = getViewportAtDepth(camera, camera.position.z);

    initParticles();

    const geometry = createGeometry(props.particleShape);
    const material = new MeshBasicMaterial({ color: props.color });
    mesh = new InstancedMesh(geometry, material, props.count);
    scene.add(mesh);

    dummy = new Object3D();
    clock = new Clock();

    window.addEventListener('pointermove', onPointerMove);
    window.addEventListener('resize', onResize);
    document.addEventListener('visibilitychange', onVisibilityChange);

    animate();
  } catch (error) {
    // WebGL initialization failed (likely due to strict performance checks or missing hardware support)
    // Log as warning and fallback to CSS
    console.warn('Antigravity: WebGL initialization failed, switching to CSS fallback.', error);
    webglAvailable.value = false;
    cleanup();
  }
}

function onPointerMove(event: PointerEvent) {
  const container = containerRef.value;
  if (!container) return;

  const rect = container.getBoundingClientRect();
  pointer.x = ((event.clientX - rect.left) / rect.width) * 2 - 1;
  pointer.y = -((event.clientY - rect.top) / rect.height) * 2 + 1;
}

function onResize() {
  const container = containerRef.value;
  if (!container || !renderer || !camera) return;

  const { clientWidth, clientHeight } = container;
  camera.aspect = clientWidth / clientHeight;
  camera.updateProjectionMatrix();
  renderer.setSize(clientWidth, clientHeight);
  
  // Update viewport size for particle redistribution
  currentViewport = getViewportAtDepth(camera, camera.position.z);
}

function onVisibilityChange() {
  isVisible = document.visibilityState === 'visible';
  if (isVisible && !animationFrameId) {
    clock.getDelta(); // Reset delta to avoid large jump
    animate();
  }
}

function animate() {
  if (!isVisible) {
    animationFrameId = 0;
    return; // Stop animation when tab is hidden
  }
  
  animationFrameId = requestAnimationFrame(animate);

  if (!mesh || !camera || !renderer || !scene) return;

  const delta = clock.getDelta();
  const elapsedTime = clock.getElapsedTime();

  const mouseDist = Math.sqrt(
    Math.pow(pointer.x - lastMousePos.x, 2) + Math.pow(pointer.y - lastMousePos.y, 2)
  );

  if (mouseDist > 0.001) {
    lastMouseMoveTime = Date.now();
    lastMousePos = { x: pointer.x, y: pointer.y };
  }

  // Use cached currentViewport
  let destX = (pointer.x * currentViewport.width) / 2;
  let destY = (pointer.y * currentViewport.height) / 2;

  if (props.autoAnimate && Date.now() - lastMouseMoveTime > 2000) {
    destX = Math.sin(elapsedTime * 0.5) * (currentViewport.width / 4);
    destY = Math.cos(elapsedTime * 0.5 * 2) * (currentViewport.height / 4);
  }

  // Frame-rate independent smoothing (faster response)
  const smoothFactor = 1 - Math.exp(-8.0 * delta);
  
  virtualMouse.x += (destX - virtualMouse.x) * smoothFactor;
  virtualMouse.y += (destY - virtualMouse.y) * smoothFactor;

  const targetX = virtualMouse.x;
  const targetY = virtualMouse.y;

  const globalRotation = elapsedTime * props.rotationSpeed;

  // Faster particle response
  const lerpDecay = props.lerpSpeed * 120;
  const particleSmoothFactor = 1 - Math.exp(-lerpDecay * delta);

  // Cache props for faster access
  const { magnetRadius, ringRadius, waveSpeed, waveAmplitude, depthFactor, fieldStrength, pulseSpeed, particleVariance, particleSize } = props;
  const viewportWidth = currentViewport.width;
  const viewportHeight = currentViewport.height;

  for (let i = 0; i < particles.length; i++) {
    const particle = particles[i];
    if (!particle) continue;
    const t = particle.t += particle.speed * 0.5;

    // Dynamically calculate mx/my based on current viewport
    const mx = particle.nx * viewportWidth;
    const my = particle.ny * viewportHeight;
    const mz = particle.mz;
    const cz = particle.cz;

    const projectionFactor = 1 - cz * 0.02; // 1 - cz / 50
    const projectedTargetX = targetX * projectionFactor;
    const projectedTargetY = targetY * projectionFactor;

    const dx = mx - projectedTargetX;
    const dy = my - projectedTargetY;
    const distSq = dx * dx + dy * dy;

    let targetPosX = mx;
    let targetPosY = my;
    let targetPosZ = mz * depthFactor;

    if (distSq < magnetRadius * magnetRadius) {
      const angle = Math.atan2(dy, dx) + globalRotation;
      const wave = Math.sin(t * waveSpeed + angle) * (0.5 * waveAmplitude);
      const deviation = particle.randomRadiusOffset * (5 / (fieldStrength + 0.1));
      const currentRingRadius = ringRadius + wave + deviation;

      targetPosX = projectedTargetX + currentRingRadius * Math.cos(angle);
      targetPosY = projectedTargetY + currentRingRadius * Math.sin(angle);
      targetPosZ = mz * depthFactor + Math.sin(t) * waveAmplitude * depthFactor;
    }

    particle.cx += (targetPosX - particle.cx) * particleSmoothFactor;
    particle.cy += (targetPosY - particle.cy) * particleSmoothFactor;
    particle.cz += (targetPosZ - particle.cz) * particleSmoothFactor;

    dummy.position.set(particle.cx, particle.cy, particle.cz);
    dummy.lookAt(projectedTargetX, projectedTargetY, particle.cz);
    dummy.rotateX(1.5707963267948966); // Math.PI / 2

    const cdx = particle.cx - projectedTargetX;
    const cdy = particle.cy - projectedTargetY;
    const currentDistToMouse = Math.sqrt(cdx * cdx + cdy * cdy);
    const distFromRing = Math.abs(currentDistToMouse - ringRadius);
    const scaleFactor = Math.max(0, Math.min(1, 1 - distFromRing * 0.1));

    const finalScale = scaleFactor * (0.8 + Math.sin(t * pulseSpeed) * 0.2 * particleVariance) * particleSize;
    dummy.scale.setScalar(finalScale);

    dummy.updateMatrix();
    mesh!.setMatrixAt(i, dummy.matrix);
  }

  mesh.instanceMatrix.needsUpdate = true;

  renderer.render(scene, camera);
}

function cleanup() {
  if (animationFrameId) {
    cancelAnimationFrame(animationFrameId);
  }

  const container = containerRef.value;
  // Cleanup all listeners
  window.removeEventListener('pointermove', onPointerMove);
  window.removeEventListener('resize', onResize);
  document.removeEventListener('visibilitychange', onVisibilityChange);

  if (mesh) {
    mesh.geometry.dispose();
    (mesh.material as Material).dispose();
  }

  if (renderer) {
    renderer.dispose();
    if (container && renderer.domElement.parentNode === container) {
      container.removeChild(renderer.domElement);
    }
  }

  renderer = null;
  scene = null;
  camera = null;
  mesh = null;
}

onMounted(setupScene);
onUnmounted(cleanup);

watch(
  () => props.color,
  (newColor: string) => {
    if (mesh) {
      (mesh.material as MeshBasicMaterial).color.set(newColor);
    }
  }
);

watch(
  () => [props.count, props.particleShape],
  () => {
    cleanup();
    setupScene();
  }
);
</script>

<template>
  <div ref="containerRef" class="antigravity-container" :style="fallbackStyle" />
</template>

<style scoped>
.antigravity-container {
  position: relative;
  width: 100%;
  height: 100%;
}
</style>
