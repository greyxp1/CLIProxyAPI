const $ = selector => document.querySelector(selector)
const names = { claude: 'Claude', codex: 'Codex' }
// T3 Code's provider marks, so both apps look the same.
const icons = {
  claude: '<svg viewBox="0 0 256 257" fill="#d97757" aria-hidden="true"><path d="m50.228 170.321 50.357-28.257.843-2.463-.843-1.361h-2.462l-8.426-.518-28.775-.778-24.952-1.037-24.175-1.296-6.092-1.297L0 125.796l.583-3.759 5.12-3.434 7.324.648 16.202 1.101 24.304 1.685 17.629 1.037 26.118 2.722h4.148l.583-1.685-1.426-1.037-1.101-1.037-25.147-17.045-27.22-18.017-14.258-10.37-7.713-5.25-3.888-4.925-1.685-10.758 7-7.713 9.397.649 2.398.648 9.527 7.323 20.35 15.75L94.817 91.9l3.889 3.24 1.555-1.102.195-.777-1.75-2.917-14.453-26.118-15.425-26.572-6.87-11.018-1.814-6.61c-.648-2.723-1.102-4.991-1.102-7.778l7.972-10.823L71.42 0 82.05 1.426l4.472 3.888 6.61 15.101 10.694 23.786 16.591 32.34 4.861 9.592 2.592 8.879.973 2.722h1.685v-1.556l1.36-18.211 2.528-22.36 2.463-28.776.843-8.1 4.018-9.722 7.971-5.25 6.222 2.981 5.12 7.324-.713 4.73-3.046 19.768-5.962 30.98-3.889 20.739h2.268l2.593-2.593 10.499-13.934 17.628-22.036 7.778-8.749 9.073-9.657 5.833-4.601h11.018l8.1 12.055-3.628 12.443-11.342 14.388-9.398 12.184-13.48 18.147-8.426 14.518.778 1.166 2.01-.194 30.46-6.481 16.462-2.982 19.637-3.37 8.88 4.148.971 4.213-3.5 8.62-20.998 5.184-24.628 4.926-36.682 8.685-.454.324.519.648 16.526 1.555 7.065.389h17.304l32.21 2.398 8.426 5.574 5.055 6.805-.843 5.184-12.962 6.611-17.498-4.148-40.83-9.721-14-3.5h-1.944v1.167l11.666 11.406 21.387 19.314 26.767 24.887 1.36 6.157-3.434 4.86-3.63-.518-23.526-17.693-9.073-7.972-20.545-17.304h-1.36v1.814l4.73 6.935 25.017 37.59 1.296 11.536-1.814 3.76-6.481 2.268-7.13-1.297-14.647-20.544-15.1-23.138-12.185-20.739-1.49.843-7.194 77.448-3.37 3.953-7.778 2.981-6.48-4.925-3.436-7.972 3.435-15.749 4.148-20.544 3.37-16.333 3.046-20.285 1.815-6.74-.13-.454-1.49.194-15.295 20.999-23.267 31.433-18.406 19.702-4.407 1.75-7.648-3.954.713-7.064 4.277-6.286 25.47-32.405 15.36-20.092 9.917-11.6-.065-1.686h-.583L44.07 198.125l-12.055 1.555-5.185-4.86.648-7.972 2.463-2.593 20.35-13.999-.064.065Z"/></svg>',
  codex: '<svg viewBox="100 100 411 411" fill="currentColor" aria-hidden="true"><path fill-rule="evenodd" d="M252.794 108.802C289.191 99.0484 326.265 110.305 351.148 135.135C385.113 126.072 422.85 134.862 449.492 161.505C476.136 188.149 484.925 225.888 475.862 259.85V259.854C500.696 284.735 511.95 321.81 502.198 358.207C492.447 394.602 464.161 421.084 430.215 430.217C421.083 464.162 394.603 492.448 358.206 502.199C321.812 511.951 284.734 500.693 259.852 475.864C225.887 484.927 188.15 476.137 161.507 449.495C134.864 422.851 126.073 385.111 135.136 351.149C110.304 326.266 99.0496 289.192 108.801 252.795C118.552 216.4 146.84 189.918 180.784 180.785C189.917 146.841 216.396 118.553 252.794 108.802ZM374.292 407.145C374.292 411.271 372.092 415.086 368.517 417.148L283.723 466.102C302.487 480.585 327.555 486.459 352.217 479.852C386.997 470.532 410.068 439.312 410.555 405.006V317.717C410.555 315.08 409.125 312.621 406.843 311.303L374.292 292.509V407.145ZM251.868 415.897C248.296 417.959 243.893 417.959 240.317 415.897L155.526 366.942C152.366 390.436 159.811 415.08 177.866 433.136H177.863C203.325 458.594 241.896 462.962 271.85 446.232L347.449 402.586C349.735 401.268 351.148 398.8 351.148 396.163V358.579L251.868 415.897ZM368.602 220.628C366.319 219.309 363.474 219.318 361.191 220.637L328.641 239.431L427.921 296.749C431.496 298.811 433.697 302.627 433.697 306.752V404.661C455.622 395.654 473.244 376.881 479.851 352.218C489.169 317.442 473.668 281.85 444.201 264.274L368.602 220.628ZM177.303 206.34C155.377 215.348 137.756 234.122 131.148 258.783C121.832 293.561 137.331 329.153 166.799 346.727L242.398 390.373C244.68 391.692 247.525 391.684 249.807 390.366L282.357 371.572L183.078 314.253C179.504 312.189 177.303 308.375 177.303 304.251V206.34ZM259.849 279.145V331.858L305.5 358.213L351.15 331.858V279.145L305.5 252.789L259.849 279.145ZM327.276 144.9C308.512 130.418 283.445 124.543 258.782 131.15C224.002 140.471 200.931 171.691 200.445 205.995V293.286C200.445 295.923 201.875 298.381 204.158 299.7L236.707 318.493V203.856C236.707 199.731 238.909 195.916 242.483 193.853L327.276 144.9ZM433.137 177.867C407.675 152.407 369.103 148.038 339.149 164.769L263.55 208.415C261.265 209.734 259.852 212.202 259.852 214.838V252.423L359.132 195.105C362.703 193.041 367.108 193.041 370.682 195.105L455.473 244.06C458.635 220.567 451.189 195.922 433.135 177.867H433.137Z"/></svg>'
}
const view = { key: sessionStorage.getItem('adminKey') || '', data: null, joining: null, legacyOpen: new Set(), modelsOpen: new Set() }
const status = $('#status')
const report = error => { status.textContent = error.message }

function element(tag, className, ...children) {
  const node = document.createElement(tag)
  if (className) node.className = className
  node.append(...children.filter(child => child != null))
  return node
}

function button(text, className, run) {
  const node = element('button', className, text)
  node.type = 'button'
  node.addEventListener('click', async () => {
    node.disabled = true
    try {
      await run()
    } catch (error) {
      report(error)
    } finally {
      node.disabled = false
    }
  })
  return node
}

async function call(path, body, token = view.key) {
  const response = await fetch(path, {
    method: body ? 'POST' : 'GET',
    headers: { ...(token && { Authorization: `Bearer ${token}` }), ...(body && { 'Content-Type': 'application/json' }) },
    body: body && JSON.stringify(body)
  })
  const data = await response.json()
  if (!response.ok) throw new Error(data.error || 'Request failed.')
  return data
}

async function load() {
  view.data = await call('/proxy/api/status')
  render()
}

const plural = (count, word) => `${count} ${word}${count === 1 ? '' : 's'}`
const compact = new Intl.NumberFormat(undefined, { notation: 'compact', maximumFractionDigits: 1 })

// Short countdown, such as 4d, 5h or 12m.
function short(time) {
  const minutes = Math.max(0, Math.round((time - Date.now()) / 60000))
  if (minutes < 60) return `${minutes}m`
  if (minutes < 1440) return `${Math.round(minutes / 60)}h`
  return `${Math.round(minutes / 1440)}d`
}

function ago(seconds) {
  if (!seconds) return 'never'
  const minutes = Math.round((Date.now() / 1000 - seconds) / 60)
  if (minutes < 1) return 'just now'
  if (minutes < 60) return `${minutes}m ago`
  if (minutes < 1440) return `${Math.round(minutes / 60)}h ago`
  return `${Math.round(minutes / 1440)}d ago`
}

const left = window => Math.max(0, 100 - window.used)
const tone = percent => percent >= 50 ? 't-good' : percent >= 20 ? 't-warn' : 't-bad'
const find = (account, name) => account.windows.find(window => window.name === name)
const resetOf = window => window?.reset ? Date.parse(window.reset) : null
const hours = { five_hour: 5, seven_day: 168 }
const planName = plan => ({ plus: 'Plus', prolite: 'Pro Lite' })[plan] || plan
// Nominal capacity estimates, not provider-reported quota amounts. Unknown plans have no estimate.
const capacityOf = account => account.provider === 'codex' ? ({ plus: 1, prolite: 5 })[account.plan] : 1
const estimated = list => new Set(list.map(capacityOf)).size > 1
const capacityHint = list => `Estimated capacity weights: ${[...new Set(list.map(account => `${planName(account.plan)} ${capacityOf(account)}×`))].join(', ')}. Actual limits can vary.`

// Share of the window still to come. Quota left beyond this mark is ahead of an even pace.
function paceOf(window) {
  const reset = resetOf(window)
  if (!reset) return null
  const length = (hours[window.name] || hours.seven_day) * 3600000
  return Math.min(Math.max((reset - Date.now()) / length, 0), 1)
}

// Golden-angle hues keep every member's colour distinct however many there are.
const colorOf = name => `hsl(${(view.data.members.findIndex(member => member.name === name) * 137.508 + 230) % 360} 70% 66%)`

function dot(name) {
  const node = element('i', 'dot')
  node.style.background = colorOf(name)
  return node
}

// A thin capsule with the value above it, like the iOS battery and Screen Time views.
function quota(window, badge) {
  if (!window) return element('div', 'quota muted', '–')
  const percent = left(window)
  const reset = resetOf(window)
  const fill = element('i', tone(percent))
  fill.style.width = `${percent}%`
  const track = element('div', 'track', fill)
  const pace = paceOf(window)
  if (pace != null) {
    const mark = element('b', 'pace')
    mark.style.left = `${pace * 100}%`
    track.append(mark)
  }
  const note = element('span', 'note muted', badge, reset ? `resets in ${short(reset)}` : 'not started')
  return element('div', 'quota', element('div', 'line', element('b', 'num', `${Math.round(percent)}%`), note), track)
}

// How much of the account's weekly limit each member used in the current window.
function shares(account) {
  const entries = Object.entries(account.members || {}).filter(([, used]) => used >= 0.5).sort((a, b) => b[1] - a[1])
  const strip = element('div', 'track')
  for (const [name, used] of entries) {
    const part = element('i')
    part.style.width = `${used}%`
    part.style.background = colorOf(name)
    part.title = `${name}: ${Math.round(used)}%`
    strip.append(part)
  }
  const used = entries.reduce((sum, [, value]) => sum + value, 0)
  return element('div', 'shares', element('div', 'line', element('b', 'num', `${Math.round(used)}%`), element('span', 'muted small', 'used')), strip)
}

function verdict(list) {
  const usable = list.filter(account => account.state !== 'disabled')
  if (!usable.length) return ['v-bad', 'No enabled accounts']
  if (!usable.some(account => account.state === 'available' || account.state === 'unknown quota')) return ['v-bad', 'Waiting for quota']
  const week = poolWindow(usable, 'seven_day')
  if (week?.pace != null) {
    if (week.left < week.pace * 100 - 10) return ['v-warn', 'Above pace']
    return ['v-ok', week.left > week.pace * 100 + 10 ? 'Below pace' : 'On pace']
  }
  return ['v-ok', 'Available']
}

function accountRow(account, owner) {
  const disabled = account.state === 'disabled'
  const email = account.email || account.id
  const label = element(owner ? 'button' : 'span', owner ? 'email editable' : 'email', account.name || email)
  label.title = owner ? `Rename account · ${email}` : email
  const title = element('div', 'name', label)
  if (account.next) title.append(element('span', 'badge next', 'Next'))
  else if (account.state !== 'available') title.append(element('span', 'badge', account.state[0].toUpperCase() + account.state.slice(1)))
  if (owner) {
    label.type = 'button'
    label.setAttribute('aria-label', `Rename account ${account.name || email}`)
    label.addEventListener('click', () => editName(account.name || '', email, title,
      name => call('/proxy/api/accounts', { id: account.id, name })))
  }
  const details = [`Used ${ago(account.last_used)}`]
  if (account.plan) details.push(planName(account.plan))
  for (const window of account.windows.filter(window => window.name.startsWith('seven_day_'))) {
    details.push(`${window.name.replace('seven_day_', '')} ${Math.round(left(window))}% left`)
  }
  if (account.retry_at) details.push(`retry in ${short(Date.parse(account.retry_at))}`)
  const identity = element('div', 'identity')
  if (owner) {
    const toggle = element('input', 'switch')
    toggle.type = 'checkbox'
    toggle.checked = !disabled
    toggle.setAttribute('role', 'switch')
    toggle.setAttribute('aria-label', `Route requests to ${email}`)
    toggle.dataset.account = account.id
    toggle.addEventListener('change', async () => {
      const enabled = toggle.checked
      toggle.disabled = true
      try {
        await call('/proxy/api/accounts', { id: account.id, disabled: !enabled })
        const focused = document.activeElement === toggle
        await load().catch(report)
        if (focused) [...document.querySelectorAll('.switch')].find(input => input.dataset.account === account.id)?.focus()
      } catch (error) {
        toggle.checked = !enabled
        report(error)
      } finally {
        toggle.disabled = false
      }
    })
    identity.append(toggle)
  }
  // Who used it in the last ten minutes, each name in their colour from the members panel.
  const usedBy = (account.last_members || []).flatMap((name, index) => {
    const node = element('span', 'by', name)
    node.style.color = colorOf(name)
    return [index ? ', ' : ' by ', node]
  })
  const [used, ...rest] = details
  identity.append(element('div', 'account', title, element('div', 'muted small', used, ...usedBy, rest.length ? ` · ${rest.join(' · ')}` : null)))
  const cells = [identity]
  if (!disabled) {
    const week = find(account, 'seven_day')
    const session = find(account, 'five_hour')
    const sessionCell = account.provider === 'codex' && week && !session
      ? element('div', 'quota no-limit muted', 'No limit')
      : quota(session)
    cells.push(quota(week, bankedResets(account, owner)), sessionCell, shares(account))
  }
  return element('div', disabled ? 'row off' : 'row', ...cells)
}

const checkIcon = 'm5 12 5 5L20 7'
const crossIcon = 'M6 6l12 12M18 6 6 18'

// Swaps the name for an input in the same spot, so nothing else moves.
function editName(value, placeholder, title, saveName, overlay = false) {
  const original = [...title.childNodes]
  const restore = () => {
    if (overlay) form.remove()
    else title.replaceChildren(...original)
    title.querySelector('.editable')?.focus()
  }
  const form = element('form', overlay ? 'name-editor overlay' : 'name-editor')
  const input = element('input')
  input.value = value
  input.placeholder = placeholder
  input.maxLength = 80
  input.setAttribute('aria-label', 'Display name')
  const save = element('button', 'icon-button confirm')
  save.type = 'submit'
  setIcon(save, 'Save name', checkIcon)
  const cancel = iconButton('Cancel', crossIcon, '', restore)
  input.addEventListener('keydown', event => { if (event.key === 'Escape') cancel.click() })
  form.append(input, save, cancel)
  form.addEventListener('submit', async event => {
    event.preventDefault()
    save.disabled = true
    try {
      await saveName(input.value)
      await load()
    } catch (error) {
      report(error)
    } finally {
      save.disabled = false
    }
  })
  if (overlay) title.append(form)
  else title.replaceChildren(form)
  input.focus()
  input.select()
}

// A badge beside the weekly reset time. The owner can use the reset that expires first.
function bankedResets(account, owner) {
  if (!account.credits) return null
  const credits = [...(account.reset_credits || [])].sort((a, b) => Date.parse(a.expires_at) - Date.parse(b.expires_at))
  const expiry = account.credit_expiry ? ` · next expires in ${short(Date.parse(account.credit_expiry))}` : ''
  const label = `${plural(account.credits, 'banked reset')}${expiry}`
  if (!owner || !credits.length) {
    const badge = element('span', 'reset-badge', `↺ ${account.credits}`)
    badge.dataset.tip = label
    badge.tabIndex = 0
    return badge
  }
  const badge = button(`↺ ${account.credits}`, 'reset-badge', async () => {
    if (!await ask('Use a banked reset?', `This restores ${account.name || account.email}'s limits now. ${label}.`, 'Use reset')) return
    const result = await call('/proxy/api/resets', { id: account.id, credit_id: credits[0].id })
    const messages = { reset: 'Reset used.', already_redeemed: 'This reset was already used.', nothing_to_reset: 'No eligible usage window needs resetting.', no_credit: 'This reset is no longer available.' }
    status.textContent = messages[result.code] || 'Reset request completed.'
    if (!result.refreshed) status.textContent += ' Limits could not refresh yet.'
    await load().catch(() => { status.textContent += ' The page could not refresh yet.' })
  })
  badge.dataset.tip = `${label} · click to use`
  return badge
}

// Use the same capacity weights for quota left, pace and member contributions.
function poolWindow(usable, name) {
  if (!usable.length || usable.some(account => !capacityOf(account) || !find(account, name))) return null
  const capacity = usable.reduce((sum, account) => sum + capacityOf(account), 0)
  return {
    left: usable.reduce((sum, account) => sum + left(find(account, name)) * capacityOf(account), 0) / capacity,
    pace: usable.reduce((sum, account) => sum + (paceOf(find(account, name)) ?? 1) * capacityOf(account), 0) / capacity,
    estimated: estimated(usable)
  }
}

function stat(window, label, empty = '–', hint = '') {
  const caption = element('span', 'stat-label muted', `${label}${window?.estimated ? ' · estimated' : ''}`)
  const node = element('div', 'stat')
  if (hint) node.title = hint
  if (!window) {
    node.className += ' empty'
    node.append(caption, element('span', 'muted', empty))
    return node
  }
  const fill = element('i', tone(window.left))
  fill.style.width = `${window.left}%`
  const mark = element('b', 'pace')
  mark.style.left = `${window.pace * 100}%`
  node.append(caption, element('b', 'num', `${Math.round(window.left)}%`), element('div', 'track', fill, mark))
  return node
}

function providerCard(provider, list, owner) {
  const usable = list.filter(account => account.state !== 'disabled')
  const hint = estimated(usable) && usable.every(capacityOf) ? capacityHint(usable) : ''
  const limited = usable.filter(account => !(account.provider === 'codex' && find(account, 'seven_day') && !find(account, 'five_hour')))
  const session = poolWindow(limited, 'five_hour')
  const plans = [...new Set(limited.map(account => account.plan))]
  const sessionLabel = limited.length < usable.length
    ? `${plans.length === 1 && plans[0] ? planName(plans[0]) : 'Capped accounts'} · 5 hours left`
    : 'Session left'
  const logo = element('span', 'logo')
  logo.innerHTML = icons[provider]
  const [level, text] = verdict(list)
  const heading = element('h2', '', names[provider])
  if (view.joining?.provider !== provider) {
    heading.append(iconButton(`Add ${names[provider]} account`, plusIcon, '', async () => {
      const result = await call('/proxy/api/join/start', { provider }, '')
      view.joining = { provider, session: result.session, url: result.url }
      render()
    }))
  }
  const models = (view.data.models || []).filter(model => model.provider === provider)
  const modelsOpen = view.modelsOpen.has(provider)
  if (owner && models.length) {
    const current = models.filter(model => !model.older)
    const toggle = iconButton(`${modelsOpen ? 'Hide' : 'Show'} ${names[provider]} models · ${current.filter(model => model.enabled).length} of ${current.length} on`,
      'M12 3 3 8l9 5 9-5-9-5ZM3 13l9 5 9-5M3 17l9 5 9-5', modelsOpen ? 'open' : '', () => {
        if (modelsOpen) view.modelsOpen.delete(provider)
        else view.modelsOpen.add(provider)
        render()
      })
    toggle.setAttribute('aria-expanded', modelsOpen)
    heading.append(toggle)
  }
  const head = element('div', 'head',
    element('div', 'identity', logo, element('div', 'title', heading, element('p', `verdict ${level}`, text))),
    stat(poolWindow(usable, 'seven_day'), 'Weekly left', '–', hint),
    stat(session, sessionLabel, usable.length && !limited.length ? 'No five-hour cap' : '–', session?.estimated ? capacityHint(limited) : ''),
    poolShare(usable))
  const controls = element('div', 'chips', ...limitChips(provider, view.data.options?.[provider], owner))
  const card = element('section', 'card', head,
    element('div', 'row labels', element('div', '', controls), element('div', '', 'Weekly'), element('div', '', 'Session'), element('div', '', 'Who used it this week')),
    ...(modelsOpen ? [modelGroup(provider, models)] : []),
    ...list.map(account => accountRow(account, owner)))
  if (view.joining?.provider === provider) card.append(joinPanel())
  return card
}

// Who used how much of a provider's pool this week. Every bar shares one scale, the whole pool,
// so members compare at a glance.
function poolShare(list) {
  if (!list.every(account => capacityOf(account) && find(account, 'seven_day'))) return element('div', 'shares muted', 'Usage shown per account')
  const capacity = (list.reduce((sum, account) => sum + capacityOf(account), 0) || 1) * 100
  const totals = {}
  for (const account of list) {
    for (const [name, used] of Object.entries(account.members || {})) totals[name] = (totals[name] || 0) + used * capacityOf(account)
  }
  const ranked = Object.entries(totals).filter(([, used]) => used >= 0.5).sort((a, b) => b[1] - a[1])
  const used = Object.values(totals).reduce((sum, value) => sum + value, 0)
  // Even pace would have used this much of the pool by now.
  const expected = 100 - poolWindow(list, 'seven_day')?.pace * 100
  const stack = element('div', 'track')
  if (Number.isFinite(expected)) {
    const mark = element('b', 'pace')
    mark.style.left = `${expected}%`
    stack.append(mark)
  }
  for (const [name, value] of ranked) {
    const part = element('i')
    part.style.width = `${value / capacity * 100}%`
    part.style.background = colorOf(name)
    part.title = `${name}: ${Math.round(value / capacity * 100)}%${estimated(list) ? ' estimated' : ''}`
    stack.append(part)
  }
  const node = element('div', 'shares pool', element('div', 'line', element('b', 'num', `${Math.round(used / capacity * 100)}%`), element('span', 'muted small', estimated(list) ? 'used · estimated' : 'used')), stack)
  if (estimated(list)) node.title = capacityHint(list)
  return node
}

function memberRow(member, owner) {
  const editable = owner && member.managed
  const name = element(editable ? 'button' : 'span', editable ? 'who editable' : 'who', member.name)
  name.title = editable ? `Rename ${member.name}` : member.name
  const title = element('div', 'member-name', name)
  // Each fact stays on one line, so a narrow panel only breaks between them.
  const facts = [member.today_total ? `${compact.format(member.today_total)} today` : null, ago(member.last_used)].filter(Boolean)
    .flatMap((fact, index) => [index ? ' · ' : null, element('span', 'fact', fact)]).filter(Boolean)
  const actions = element('div', 'member-actions')
  const line = element('div', 'mhead', dot(member.name), title, actions,
    element('div', 'usage', element('b', 'num', compact.format(member.total)), element('div', 'muted small', ...facts)))
  const node = element('div', 'member', line)
  if (editable) {
    name.type = 'button'
    name.setAttribute('aria-label', `Rename member ${member.name}`)
    name.addEventListener('click', () => editName(member.name, member.name, line,
      newName => call('/proxy/api/clients', { action: 'rename', name: member.name, new_name: newName }), true))
    actions.append(
      iconButton(`Show key for ${member.name}`, 'M15.75 5.25a3 3 0 0 1 3 3m3 0a6 6 0 0 1-7.029 5.912c-.563-.097-1.159.026-1.563.43L10.5 16.5H8.25v2.25H6V21H3v-3l6.908-6.908c.404-.404.527-1 .43-1.563A6 6 0 1 1 21.75 8.25Z', '',
        async () => showPrompt(member.name, (await call('/proxy/api/clients', { action: 'show', name: member.name })).prompt)),
      iconButton(`Revoke key for ${member.name}`, 'M4 7h16M10 3h4l2 4M6 7l1 14h10l1-14M10 11v6M14 11v6', 'danger', async () => {
        if (!await ask(`Revoke ${member.name}'s key?`, 'Their setup stops working immediately.', 'Revoke', true)) return
        await call('/proxy/api/clients', { action: 'revoke', name: member.name })
        await load()
      }))
  }
  return node
}

const plusIcon = 'M12 5v14M5 12h14'

function iconButton(label, path, className, run) {
  const node = button('', `icon-button ${className}`, run)
  setIcon(node, label, path)
  return node
}

function setIcon(node, label, path) {
  node.title = label
  node.setAttribute('aria-label', label)
  node.innerHTML = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="${path}"/></svg>`
}

// Posts an owner control's new state, reloads and keeps keyboard focus on the same control.
async function saveControl(control, path, body, revert) {
  control.disabled = true
  try {
    await call(path, body)
    const focused = document.activeElement === control
    await load().catch(report)
    if (focused) document.querySelector(`[data-key="${CSS.escape(control.dataset.key)}"]`)?.focus()
  } catch (error) {
    revert()
    report(error)
  } finally {
    control.disabled = false
  }
}

function ownerSwitch(label, checked, key, path, body) {
  const toggle = element('input', 'switch')
  toggle.type = 'checkbox'
  toggle.checked = checked
  toggle.setAttribute('role', 'switch')
  toggle.setAttribute('aria-label', label)
  toggle.dataset.key = key
  toggle.addEventListener('change', () =>
    saveControl(toggle, path, body(toggle.checked), () => { toggle.checked = !toggle.checked }))
  return toggle
}

// The owner's switch for one model. Off blocks it for everyone and hides it from the model list.
function modelRow(model) {
  const label = element('span', 'model-name', model.name || model.id)
  label.title = model.older ? `${model.id} · a newer version replaces it` : model.id
  return element('label', model.enabled ? 'model' : 'model off',
    ownerSwitch(`Allow ${model.id}`, model.enabled, `model:${model.id}`, '/proxy/api/models', enabled => ({ id: model.id, enabled })), label)
}

const effortNames = { minimal: 'Minimal', low: 'Low', medium: 'Medium', high: 'High', xhigh: 'XHigh', max: 'Max' }

// The effort levels the provider's models offer, plus fast mode. A request at a switched-off level
// runs at the nearest allowed one below. Members see the same chips without being able to change them.
function limitChips(provider, options, owner) {
  if (!options) return []
  const chip = (setting, text, enabled) => {
    const node = element(owner ? 'button' : 'span', enabled ? 'chip on' : 'chip', text)
    node.title = `${text} ${enabled ? 'allowed' : 'off'}`
    if (!owner) return node
    node.type = 'button'
    node.dataset.key = `option:${provider}:${setting}`
    node.setAttribute('aria-pressed', enabled)
    node.addEventListener('click', () => {
      node.classList.toggle('on')
      saveControl(node, '/proxy/api/options', { provider, setting, enabled: !enabled }, () => node.classList.toggle('on'))
    })
    return node
  }
  return [...options.efforts.map(({ level, enabled }) => chip(level, effortNames[level] || level, enabled)), chip('fast', 'Fast', options.fast)]
}

// Current models, then older versions behind a toggle, as T3 Code does. The toggle sits above the
// older models, so opening them never moves it.
function modelGroup(provider, models) {
  const panel = element('div', 'model-panel', element('div', 'model-grid', ...models.filter(model => !model.older).map(modelRow)))
  const legacy = models.filter(model => model.older)
  if (!legacy.length) return panel
  const open = view.legacyOpen.has(provider)
  const on = legacy.filter(model => model.enabled).length
  const toggle = button(open ? 'Hide legacy models' : `Show ${plural(legacy.length, 'legacy model')}${on ? ` · ${on} on` : ''}`, 'text legacy-toggle', () => {
    if (open) view.legacyOpen.delete(provider)
    else view.legacyOpen.add(provider)
    render()
  })
  toggle.setAttribute('aria-expanded', open)
  panel.append(toggle, ...(open ? [element('div', 'model-grid legacy-grid', ...legacy.map(modelRow))] : []))
  return panel
}

function render() {
  const { accounts, members, role } = view.data
  const owner = role === 'owner'
  const checked = Math.max(...accounts.map(account => Date.parse(account.checked_at || 0)))
  $('#updated').textContent = checked > 0 ? `Limits checked ${ago(checked / 1000)}` : ''
  const byProvider = provider => accounts.filter(account => account.provider === provider)
  $('#providers').replaceChildren(...['claude', 'codex'].map(provider => providerCard(provider, byProvider(provider), owner)))
  const sorted = [...members].sort((a, b) => b.total - a.total || (b.last_used || 0) - (a.last_used || 0))
  $('#members').replaceChildren(...sorted.map(member => memberRow(member, owner)))
  if (!owner) $('#add').hidden = true
  $('#add-member').hidden = !owner || !$('#add').hidden
  setIcon($('#admin'), owner ? 'Sign out' : 'Admin', owner
    ? 'M9 5H5v14h4m5-14 7 7-7 7M21 12H9'
    : 'M20 21v-2a7 7 0 0 0-14 0v2M12 3a4 4 0 1 0 0 8 4 4 0 0 0 0-8Z')
}

// Sign-in runs on the server, so it works without a key or port forwarding.
function joinPanel() {
  const joining = view.joining
  const link = element('a', '', `Open ${names[joining.provider]} sign-in ↗`)
  link.href = joining.url
  link.target = '_blank'
  link.rel = 'noopener noreferrer'
  const field = element('input')
  field.spellcheck = false
  field.autocomplete = 'off'
  field.placeholder = 'Final page address or code'
  field.setAttribute('aria-label', 'Callback URL or authorization code')
  const finish = button('Finish connecting', 'primary', async () => {
    await call('/proxy/api/join/complete', { callback: field.value }, joining.session)
    status.textContent = 'Completing sign-in…'
    for (const deadline = Date.now() + 60000; Date.now() < deadline;) {
      await new Promise(resolve => setTimeout(resolve, 2000))
      const result = await call('/proxy/api/join/status', {}, joining.session)
      if (result.status === 'error') throw new Error(result.error)
      if (result.status === 'ok') {
        view.joining = null
        status.textContent = `Connected. The ${names[joining.provider]} account is now part of the pool. Its email shows on this page, and if paid extra usage is on, shared traffic may incur charges.`
        await load()
        return
      }
    }
    throw new Error('Sign-in is taking longer than expected. Start again if it does not complete.')
  })
  const cancel = button('Cancel', '', async () => { view.joining = null; render() })
  field.addEventListener('keydown', event => { if (event.key === 'Enter') finish.click() })
  return element('div', 'join', link, element('p', 'muted', 'After signing in, paste the full address of the final page here, even if it fails to load. If Claude shows a code instead, paste the code. Finish within five minutes.'),
    element('div', 'inline', field, finish, cancel))
}

// A styled stand-in for confirm(). Resolves true only when the action button is pressed.
function ask(title, text, action, danger = false) {
  const dialog = $('#ask')
  $('#ask-title').textContent = title
  $('#ask-text').textContent = text
  $('#ask-ok').textContent = action
  $('#ask-ok').className = danger ? 'danger-fill' : 'primary'
  dialog.returnValue = ''
  dialog.showModal()
  $('#ask-cancel').focus()
  return new Promise(resolve => dialog.addEventListener('close', () => resolve(dialog.returnValue === 'ok'), { once: true }))
}

function showPrompt(name, text) {
  $('#prompt-title').textContent = `Setup for ${name}`
  $('#prompt').value = text
  if (!$('#created').open) $('#created').showModal()
}

$('#admin').addEventListener('click', () => {
  if (view.data?.role === 'owner') {
    sessionStorage.removeItem('adminKey')
    view.key = ''
    $('#created').close()
    load().catch(report)
    return
  }
  $('#login').hidden = !$('#login').hidden
  if (!$('#login').hidden) $('#key').focus()
})
$('#login').addEventListener('submit', async event => {
  event.preventDefault()
  view.key = $('#key').value.trim()
  try {
    await load()
    sessionStorage.setItem('adminKey', view.key)
    $('#login').hidden = true
    $('#key').value = ''
    status.textContent = ''
  } catch (error) {
    view.key = ''
    report(error)
  }
})
// Maintenance checks limits every five minutes. Refresh checks each account now, at most once a minute.
$('#refresh').addEventListener('click', async () => {
  const node = $('#refresh')
  node.disabled = true
  await Promise.allSettled((view.data?.accounts || []).filter(account => account.state !== 'disabled').map(account => call('/proxy/api/refresh', { id: account.id }, '')))
  node.disabled = false
  load().catch(report)
})
$('#add').addEventListener('submit', async event => {
  event.preventDefault()
  const name = $('#name').value.trim()
  try {
    const data = await call('/proxy/api/clients', { action: 'add', name })
    $('#name').value = ''
    $('#add').hidden = true
    showPrompt(name, data.prompt)
    await load()
  } catch (error) {
    report(error)
  }
})
$('#add-member').addEventListener('click', () => {
  $('#add').hidden = false
  $('#add-member').hidden = true
  $('#name').focus()
})
setIcon($('#add-member'), 'Add member', plusIcon)
setIcon($('#confirm-add'), 'Add member', checkIcon)
setIcon($('#cancel-add'), 'Cancel', crossIcon)
$('#name').addEventListener('keydown', event => { if (event.key === 'Escape') $('#cancel-add').click() })
$('#cancel-add').addEventListener('click', () => {
  $('#add').hidden = true
  $('#add-member').hidden = false
  $('#name').value = ''
  $('#add-member').focus()
})
const copyIcon = 'M8 8V5a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2h-3M5 8h9a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-9a2 2 0 0 1 2-2Z'
setIcon($('#copy'), 'Copy prompt', copyIcon)
$('#copy').addEventListener('click', async () => {
  await navigator.clipboard.writeText($('#prompt').value)
  setIcon($('#copy'), 'Copied', checkIcon)
  setTimeout(() => setIcon($('#copy'), 'Copy prompt', copyIcon), 1500)
})
$('#created').addEventListener('close', () => {
  $('#prompt').value = ''
})
for (const dialog of document.querySelectorAll('dialog')) {
  dialog.addEventListener('click', event => {
    if (event.target !== dialog) return
    const bounds = dialog.getBoundingClientRect()
    if (event.clientX < bounds.left || event.clientX > bounds.right || event.clientY < bounds.top || event.clientY > bounds.bottom) dialog.close()
  })
}
$('#ask-ok').addEventListener('click', () => $('#ask').close('ok'))
$('#ask-cancel').addEventListener('click', () => $('#ask').close())
setIcon($('#refresh'), 'Refresh limits', 'M20 7v5h-5M4 17v-5h5M6 7a7 7 0 0 1 12-1l2 2M4 16l2 2a7 7 0 0 0 12-1')

// Reloading re-sorts each provider's accounts, so priority changes move rows up or down.
// A sign-in in progress is kept, since rendering rebuilds the join panel from view.joining.
setInterval(() => {
  if (!document.hidden && !view.joining) load().catch(report)
}, 60000)

load().catch(error => {
  // A stale admin key falls back to the public view.
  if (!view.key) return report(error)
  sessionStorage.removeItem('adminKey')
  view.key = ''
  load().catch(report)
})
