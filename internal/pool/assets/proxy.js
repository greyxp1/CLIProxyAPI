const $ = selector => document.querySelector(selector)
const names = { claude: 'Claude', codex: 'Codex' }
// T3 Code's provider marks, so both apps look the same.
const icons = {
  claude: '<svg viewBox="0 0 256 257" fill="#d97757" aria-hidden="true"><path d="m50.228 170.321 50.357-28.257.843-2.463-.843-1.361h-2.462l-8.426-.518-28.775-.778-24.952-1.037-24.175-1.296-6.092-1.297L0 125.796l.583-3.759 5.12-3.434 7.324.648 16.202 1.101 24.304 1.685 17.629 1.037 26.118 2.722h4.148l.583-1.685-1.426-1.037-1.101-1.037-25.147-17.045-27.22-18.017-14.258-10.37-7.713-5.25-3.888-4.925-1.685-10.758 7-7.713 9.397.649 2.398.648 9.527 7.323 20.35 15.75L94.817 91.9l3.889 3.24 1.555-1.102.195-.777-1.75-2.917-14.453-26.118-15.425-26.572-6.87-11.018-1.814-6.61c-.648-2.723-1.102-4.991-1.102-7.778l7.972-10.823L71.42 0 82.05 1.426l4.472 3.888 6.61 15.101 10.694 23.786 16.591 32.34 4.861 9.592 2.592 8.879.973 2.722h1.685v-1.556l1.36-18.211 2.528-22.36 2.463-28.776.843-8.1 4.018-9.722 7.971-5.25 6.222 2.981 5.12 7.324-.713 4.73-3.046 19.768-5.962 30.98-3.889 20.739h2.268l2.593-2.593 10.499-13.934 17.628-22.036 7.778-8.749 9.073-9.657 5.833-4.601h11.018l8.1 12.055-3.628 12.443-11.342 14.388-9.398 12.184-13.48 18.147-8.426 14.518.778 1.166 2.01-.194 30.46-6.481 16.462-2.982 19.637-3.37 8.88 4.148.971 4.213-3.5 8.62-20.998 5.184-24.628 4.926-36.682 8.685-.454.324.519.648 16.526 1.555 7.065.389h17.304l32.21 2.398 8.426 5.574 5.055 6.805-.843 5.184-12.962 6.611-17.498-4.148-40.83-9.721-14-3.5h-1.944v1.167l11.666 11.406 21.387 19.314 26.767 24.887 1.36 6.157-3.434 4.86-3.63-.518-23.526-17.693-9.073-7.972-20.545-17.304h-1.36v1.814l4.73 6.935 25.017 37.59 1.296 11.536-1.814 3.76-6.481 2.268-7.13-1.297-14.647-20.544-15.1-23.138-12.185-20.739-1.49.843-7.194 77.448-3.37 3.953-7.778 2.981-6.48-4.925-3.436-7.972 3.435-15.749 4.148-20.544 3.37-16.333 3.046-20.285 1.815-6.74-.13-.454-1.49.194-15.295 20.999-23.267 31.433-18.406 19.702-4.407 1.75-7.648-3.954.713-7.064 4.277-6.286 25.47-32.405 15.36-20.092 9.917-11.6-.065-1.686h-.583L44.07 198.125l-12.055 1.555-5.185-4.86.648-7.972 2.463-2.593 20.35-13.999-.064.065Z"/></svg>',
  codex: '<svg viewBox="100 100 411 411" fill="currentColor" aria-hidden="true"><path fill-rule="evenodd" d="M252.794 108.802C289.191 99.0484 326.265 110.305 351.148 135.135C385.113 126.072 422.85 134.862 449.492 161.505C476.136 188.149 484.925 225.888 475.862 259.85V259.854C500.696 284.735 511.95 321.81 502.198 358.207C492.447 394.602 464.161 421.084 430.215 430.217C421.083 464.162 394.603 492.448 358.206 502.199C321.812 511.951 284.734 500.693 259.852 475.864C225.887 484.927 188.15 476.137 161.507 449.495C134.864 422.851 126.073 385.111 135.136 351.149C110.304 326.266 99.0496 289.192 108.801 252.795C118.552 216.4 146.84 189.918 180.784 180.785C189.917 146.841 216.396 118.553 252.794 108.802ZM374.292 407.145C374.292 411.271 372.092 415.086 368.517 417.148L283.723 466.102C302.487 480.585 327.555 486.459 352.217 479.852C386.997 470.532 410.068 439.312 410.555 405.006V317.717C410.555 315.08 409.125 312.621 406.843 311.303L374.292 292.509V407.145ZM251.868 415.897C248.296 417.959 243.893 417.959 240.317 415.897L155.526 366.942C152.366 390.436 159.811 415.08 177.866 433.136H177.863C203.325 458.594 241.896 462.962 271.85 446.232L347.449 402.586C349.735 401.268 351.148 398.8 351.148 396.163V358.579L251.868 415.897ZM368.602 220.628C366.319 219.309 363.474 219.318 361.191 220.637L328.641 239.431L427.921 296.749C431.496 298.811 433.697 302.627 433.697 306.752V404.661C455.622 395.654 473.244 376.881 479.851 352.218C489.169 317.442 473.668 281.85 444.201 264.274L368.602 220.628ZM177.303 206.34C155.377 215.348 137.756 234.122 131.148 258.783C121.832 293.561 137.331 329.153 166.799 346.727L242.398 390.373C244.68 391.692 247.525 391.684 249.807 390.366L282.357 371.572L183.078 314.253C179.504 312.189 177.303 308.375 177.303 304.251V206.34ZM259.849 279.145V331.858L305.5 358.213L351.15 331.858V279.145L305.5 252.789L259.849 279.145ZM327.276 144.9C308.512 130.418 283.445 124.543 258.782 131.15C224.002 140.471 200.931 171.691 200.445 205.995V293.286C200.445 295.923 201.875 298.381 204.158 299.7L236.707 318.493V203.856C236.707 199.731 238.909 195.916 242.483 193.853L327.276 144.9ZM433.137 177.867C407.675 152.407 369.103 148.038 339.149 164.769L263.55 208.415C261.265 209.734 259.852 212.202 259.852 214.838V252.423L359.132 195.105C362.703 193.041 367.108 193.041 370.682 195.105L455.473 244.06C458.635 220.567 451.189 195.922 433.135 177.867H433.137Z"/></svg>'
}
const view = { key: sessionStorage.getItem('adminKey') || '', data: null, joining: null, legacyOpen: new Set(), settingsFor: null, chartProvider: 'all', forecasting: localStorage.getItem('limitsView') === 'forecast', history: null, historyAt: -Infinity }
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
  if (!response.ok) throw Object.assign(new Error(data.error || 'Request failed.'), { status: response.status })
  return data
}

// Overlapping loads render in the order they started, so a slow stale reply never replaces a newer one.
let loads = 0
async function load() {
  const ticket = ++loads
  // Weekly history for the forecasts. Maintenance records a reading every five minutes, so it is fetched less often
  // than the status, and a failed fetch keeps the last one.
  const due = performance.now() - view.historyAt > 240000
  const [data, history] = await Promise.all([call('/proxy/api/status'), due ? call('/proxy/api/history', null, '').catch(() => null) : null])
  if (ticket !== loads) return
  if (history) [view.history, view.historyAt] = [history.accounts, performance.now()]
  view.data = data
  view.receivedAt = performance.now()
  render()
}

// Advance the server's time with elapsed time, independent of the browser's clock.
const now = () => view.data.server_time + performance.now() - view.receivedAt
const plural = (count, word) => `${count} ${word}${count === 1 ? '' : 's'}`
const compact = new Intl.NumberFormat('en-US', { notation: 'compact', maximumFractionDigits: 1 })

// Short countdown, such as 4d, 5h or 12m.
function short(time) {
  const minutes = Math.max(0, Math.round((time - now()) / 60000))
  if (minutes < 60) return `${minutes}m`
  // Minutes under a day, so a window opened moments ago doesn't read as untouched.
  if (minutes < 1440) return `${Math.floor(minutes / 60)}h${minutes % 60 ? ` ${minutes % 60}m` : ''}`
  return `${Math.round(minutes / 1440)}d`
}

function ago(seconds) {
  if (!seconds) return 'never'
  const minutes = Math.round((now() / 1000 - seconds) / 60)
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
// Weekly limits of one model family, such as seven_day_fable. The server only sends those for models that can be picked.
const modelWindows = list => [...new Set(list.flatMap(account => account.windows.map(window => window.name).filter(name => name.startsWith('seven_day_'))))].sort()
const familyOf = name => name[10].toUpperCase() + name.slice(11)
// Nominal capacity estimates, not provider-reported quota amounts. Unknown plans have no estimate.
const capacityOf = account => account.provider === 'codex' ? ({ plus: 1, prolite: 5 })[account.plan] : 1

// Every enabled account counts toward the pool, exhausted ones at 0%, so the bars and the forecast agree.
const poolAccounts = list => list.filter(account => account.state !== 'disabled')

// Share of the window still to come. Quota left beyond this mark is ahead of an even pace.
function paceOf(window) {
  const reset = resetOf(window)
  if (!reset) return null
  const length = (hours[window.name] || hours.seven_day) * 3600000
  return Math.min(Math.max((reset - now()) / length, 0), 1)
}

const palette = ['#3987e5', '#199e70', '#d95926', '#9085e9', '#d55181']

// Slots follow name order among members who have used the pool, so colours hold while usage shifts.
function colorOf(name) {
  const used = view.data.members.filter(member => member.total).map(member => member.name).sort()
  const index = used.indexOf(name)
  return index < 0 ? 'var(--muted)' : palette[index % palette.length]
}

// A member's name in their colour, matching their bar parts and chart line.
function named(name, className = 'who-tag') {
  const node = element('span', className, name)
  node.style.color = colorOf(name)
  return node
}

// Quota left fills from the left in its status colour and the used part splits into the members who
// used it, so one bar answers both how much is left and who spent the rest.
function budget(percent, shares, pace) {
  const track = element('div', 'track')
  if (percent > 0) {
    const fill = element('i', percent >= 100 ? `${tone(percent)} end` : tone(percent))
    fill.style.width = `${percent}%`
    track.append(fill)
  }
  const used = 100 - percent
  const total = shares.reduce((sum, [, value]) => sum + value, 0)
  // Shares within half a point of the used part fill it, so rounding never leaves a sliver of empty track.
  const fills = total > 0 && total > used - 0.5
  const scale = fills ? used / total : 1
  // Hovering the used group widens its tiny members until they are visible.
  const group = element('span', 'used')
  group.style.setProperty('--rest', `${used - total * scale}%`)
  for (const [name, value] of shares) {
    const part = element('i')
    part.style.setProperty('--width', value * scale)
    part.style.setProperty('--c', colorOf(name))
    part.dataset.pct = value < 0.1 ? '<0.1' : value < 1 ? value.toFixed(1) : Math.round(value)
    part.dataset.member = name
    part.tabIndex = 0
    part.setAttribute('aria-label', `${name}, ${part.dataset.pct}%`)
    group.append(part)
  }
  if (fills) group.lastChild.classList.add('end')
  if (total > 0) track.append(group)
  const bar = element('div', 'bar', track)
  if (pace != null) {
    const mark = element('b', 'pace')
    mark.style.left = `${pace * 100}%`
    bar.append(mark)
  }
  return bar
}

const ranked = entries => entries.filter(([, used]) => used > 0).sort((a, b) => b[1] - a[1])

// Highlights one member's share everywhere, or clears it with null.
function focusMember(selected) {
  for (const track of document.querySelectorAll('.bar .track, #chart .plot')) highlight(track, selected)
  $('.board').classList.toggle('member-focus', !!selected)
  for (const member of document.querySelectorAll('#members .member')) member.classList.toggle('selected', member.dataset.person === selected)
}

function memberHover(node, name) {
  node.addEventListener('pointerenter', () => focusMember(name))
  node.addEventListener('pointerleave', () => focusMember(null))
  node.addEventListener('focusin', () => focusMember(name))
  node.addEventListener('focusout', () => focusMember(null))
  return node
}


function quota(window, shares = [], badge = null, pace = paceOf(window)) {
  if (!window) return element('div', 'quota muted', '–')
  const percent = left(window)
  const reset = resetOf(window)
  return element('div', 'quota',
    element('div', 'line', element('b', 'num', `${Math.round(percent)}%`), element('span', 'note muted', badge, reset ? element('span', null, element('span', 'reset-word', 'resets in '), short(reset)) : 'not started')),
    budget(percent, shares, pace))
}


function accountRow(account, owner, models) {
  const disabled = account.state === 'disabled'
  const email = account.email || account.id
  const label = element(owner ? 'button' : 'span', owner ? 'email editable' : 'email', account.name || email)
  const title = element('div', 'name', label)
  // Only plans that count for more of the pool than the base plan get a chip.
  const weight = capacityOf(account)
  if (weight > 1) title.append(element('span', 'weight', `${weight}×`))
  if (account.next) title.append(element('span', 'badge next', 'Next'))
  else if (account.state !== 'available') title.append(element('span', 'badge', account.state[0].toUpperCase() + account.state.slice(1)))
  // Billing dates are estimates. Members see them in the week before they matter, the owner always.
  const renews = Date.parse(account.renews)
  const soon = renews - now() < 7 * 86400000
  if (soon || owner && renews) {
    const verb = account.ends ? 'Ends' : 'Renews'
    const badge = element('span', `badge renewal${soon ? ' soon' : ''}${account.ends ? ' ends' : ''}`,
      element('span', 'renewal-verb', soon ? `${verb} in` : verb),
      soon ? short(renews) : new Date(renews).toLocaleDateString(undefined, { month: 'short', day: 'numeric' }))
    badge.insertAdjacentHTML('afterbegin', `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" aria-hidden="true"><path d="${calendarIcon}"/></svg>`)
    title.append(badge)
  }
  if (owner) {
    label.type = 'button'
    label.setAttribute('aria-label', `Rename account ${account.name || email}`)
    label.addEventListener('click', () => editName(account.name || '', email, title,
      name => call('/proxy/api/accounts', { id: account.id, name })))
    title.append(element('div', 'row-actions', iconButton(`Remove ${email}`, trashIcon, 'danger', async () => {
      if (!await ask(`Remove ${email}?`, 'Its sign-in is deleted from the pool. Adding it back needs a new sign-in.', 'Remove', true)) return
      await call('/proxy/api/accounts', { id: account.id, delete: true })
      await load()
    })))
  }
  const recent = account.last_used ? account.last_members || [] : []
  const details = [element('span', 'used-by', `Used ${ago(account.last_used)}`,
    ...recent.flatMap((name, index) => [index ? ', ' : ' by ', named(name)]))]
  // Maintenance checks every five minutes, so an older reading means checks are failing.
  const checked = Date.parse(account.checked_at)
  if (!disabled && now() - checked > 20 * 60000) {
    const stale = element('span', 'stale', `Limits from ${ago(checked / 1000)}`)
    stale.dataset.tip = 'Limit checks for this account are failing, so these numbers may be out of date'
    stale.tabIndex = 0
    details.push(stale)
  }
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
      // Read focus before disabling, which moves focus off the control.
      const focused = document.activeElement === toggle
      toggle.disabled = true
      try {
        await call('/proxy/api/accounts', { id: account.id, disabled: !enabled })
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
  identity.append(element('div', 'account', title, element('div', 'account-details', ...details.map(detail => element('span', '', detail)))))
  const cells = [identity]
  if (!disabled) {
    const week = find(account, 'seven_day')
    const session = find(account, 'five_hour')
    const shares = Object.entries(account.members || {}).filter(([, used]) => used > 0).sort((a, b) => b[1] - a[1])
    // Without weekly quota the session has nothing to spend until the week resets.
    const sessionCell = account.provider === 'codex' && week && !session ? element('div', 'quota')
      : week && session && left(week) === 0 ? quota({ ...session, used: 100, reset: week.reset }, [], null, null)
      : quota(session)
    // Weekly quota that resets within a day, with more left than an even pace would leave, is lost unless used now.
    const untilReset = week ? resetOf(week) - now() : 0
    const expiring = untilReset > 0 && untilReset < 86400000 && left(week) >= Math.max(10, paceOf(week) * 100)
    const weekCell = element('div', expiring ? 'account-week expiring' : 'account-week', quota(week, shares, bankedResets(account, owner)))
    if (expiring) weekCell.setAttribute('aria-description', 'Weekly quota resets soon with plenty left. Use it now.')
    sessionCell.dataset.label = 'Session'
    cells.push(weekCell, ...models.map(name => {
      const cell = find(account, name) ? quota(find(account, name)) : element('div', 'quota')
      cell.dataset.label = familyOf(name)
      return cell
    }), sessionCell)
  }
  return element('div', disabled ? 'row off' : 'row', ...cells)
}

const checkIcon = 'm5 12 5 5L20 7'

// Swaps the name for an input in the same spot, so nothing else moves.
function editName(value, placeholder, title, saveName, overlay = false) {
  const original = [...title.childNodes]
  // Enter saves and Escape or leaving the field cancels, so the editor needs no buttons.
  let done = false
  const restore = () => {
    if (done) return
    done = true
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
  let saving = false
  input.addEventListener('keydown', event => {
    if (event.key === 'Escape') { event.preventDefault(); restore() }
  })
  input.addEventListener('blur', () => { if (!saving) restore() })
  form.append(input)
  form.addEventListener('submit', async event => {
    event.preventDefault()
    saving = true
    try {
      await saveName(input.value)
      await load()
    } catch (error) {
      report(error)
      saving = false
      input.focus()
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
  const expiry = account.credit_expiry ? `Next expires in ${short(Date.parse(account.credit_expiry))}` : ''
  const label = `${plural(account.credits, 'banked reset')}${expiry ? `, ${expiry.toLowerCase()}` : ''}`
  // Every banked reset by when it expires, soonest first.
  const tip = () => [
    ...(credits.length ? credits.map(credit => {
      const at = Date.parse(credit.expires_at)
      const when = new Date(at).toLocaleString(undefined, { month: 'short', day: 'numeric', hour: 'numeric', minute: '2-digit' })
      const row = element('div', 'tip-row', `Expires ${when}`, element('span', 'v', `in ${short(at)}`))
      // Claude names its resets, which differ in what they clear. Names are long, so they wrap above the expiry.
      return credit.label ? element('div', 'tip-reset', element('div', '', credit.label), row) : row
    }) : expiry ? [element('div', 'muted', expiry)] : []),
    owner && credits.length ? element('div', 'tip-row total', 'Click to use the soonest') : null].filter(Boolean)
  if (!owner || !credits.length) {
    const badge = element('span', 'reset-badge', `↺ ${account.credits}`)
    badge.dataset.tip = ''
    badge.tip = tip
    badge.setAttribute('aria-label', label)
    badge.tabIndex = 0
    return badge
  }
  const badge = button(`↺ ${account.credits}`, 'reset-badge', async () => {
    const what = credits[0].label ? `“${credits[0].label}”` : 'the soonest banked reset'
    if (!await ask('Use a banked reset?', `This spends ${what} on ${account.name || account.email} now. ${label}.`, 'Use reset')) return
    const result = await call('/proxy/api/resets', { id: account.id, credit_id: credits[0].id })
    const messages = { reset: 'Reset used.', already_redeemed: 'This reset was already used.', nothing_to_reset: 'No eligible usage window needs resetting.', no_credit: 'This reset is no longer available.' }
    status.textContent = messages[result.code] || 'Reset request completed.'
    if (!result.refreshed) status.textContent += ' Limits could not refresh yet.'
    await load().catch(() => { status.textContent += ' The page could not refresh yet.' })
  })
  badge.dataset.tip = ''
  badge.tip = tip
  badge.setAttribute('aria-label', `Use a banked reset for ${account.name || account.email}`)
  return badge
}

// Use the same capacity weights for quota left, pace and member contributions.
function poolWindow(usable, name) {
  if (!usable.length || usable.some(account => !capacityOf(account) || !find(account, name))) return null
  const capacity = usable.reduce((sum, account) => sum + capacityOf(account), 0)
  return {
    left: usable.reduce((sum, account) => sum + left(find(account, name)) * capacityOf(account), 0) / capacity,
    pace: usable.reduce((sum, account) => sum + (paceOf(find(account, name)) ?? 1) * capacityOf(account), 0) / capacity
  }
}

// Member use of the remaining accounts, on the same capacity scale as the pool's quota left.
function poolShares(list) {
  if (!list.every(account => capacityOf(account))) return []
  const capacity = list.reduce((sum, account) => sum + capacityOf(account), 0)
  const totals = {}
  for (const account of list) for (const [name, used] of Object.entries(account.members || {})) totals[name] = (totals[name] || 0) + used * capacityOf(account) / capacity
  return ranked(Object.entries(totals))
}

function poolCell(label, pool, shares) {
  if (!pool) return element('div', 'pool empty muted', label === 'weekly' ? 'Weekly pool not reported' : '')
  const caption = element('span', 'metric-label', { weekly: 'Weekly pool', session: '5-hour session' }[label] || label)
  const number = element('div', 'pool-number', element('b', '', `${Math.round(pool.left)}`), element('span', 'percent', '%'), element('span', 'left-label', 'left'))
  return element('div', `pool ${label === 'weekly' ? 'week' : 'session'}-pool`,
    element('div', 'pool-value', number, caption), budget(pool.left, shares, pool.pace))
}

const hour = 3600000
const week = 168 * hour

// A weekday and hour within the coming week, a date further out.
function when(time) {
  return new Date(time).toLocaleString(undefined, Math.abs(time - now()) < 6 * 86400000 ? { weekday: 'short', hour: 'numeric' } : { month: 'short', day: 'numeric' })
}

// One provider's weekly pool over time, as a share of every enabled account's weekly limit. It is recorded from
// the history, then projected at the last day's burn rate. The projection spends accounts in order of their weekly
// reset, soonest first, since quota left at a reset is lost, and refills each one when its week ends. Five-hour limits
// and banked resets are left out.
function forecast(list, from, to) {
  const accounts = list.filter(account => account.state !== 'disabled' && find(account, 'seven_day'))
  if (!accounts.length || accounts.some(account => !capacityOf(account))) return null
  const start = now()
  const step = hour / 2
  // Plan weights scale each account's share, so a 5× plan counts five times as much.
  const capacity = accounts.reduce((sum, account) => sum + capacityOf(account), 0)
  // Each account's readings in time order, ending with the live one.
  const readings = accounts.map(account => {
    const window = find(account, 'seven_day')
    return [...(view.history?.[account.id] || []).filter(([at]) => at < start).map(([at, used, reset]) => ({ at, used, reset })),
      { at: start, used: window.used, reset: resetOf(window) }]
  })
  const leftAt = (reading, time) => reading.reset && reading.reset <= time ? 100 : Math.max(0, 100 - reading.used)

  // A reading covers the hour after it, so a gap in the history breaks the line instead of bridging it.
  const recorded = []
  const cursors = readings.map(() => -1)
  for (let time = Math.ceil(from / step) * step; ; time += step) {
    if (time >= start) time = start
    let total = 0
    for (const [index, list] of readings.entries()) {
      while (cursors[index] + 1 < list.length && list[cursors[index] + 1].at <= time) cursors[index]++
      const reading = list[cursors[index]]
      if (!reading || time - reading.at > hour) { total = null; break }
      total += leftAt(reading, time) * capacityOf(accounts[index])
    }
    recorded.push({ time, value: total })
    if (time === start) break
  }

  // Quota spent over the last day. A later reset time means a new week, which began with nothing used.
  const since = Math.max(start - 24 * hour, ...readings.map(list => list[0].at))
  let drained = 0
  for (const [index, list] of readings.entries()) {
    const inside = list.filter(reading => reading.at >= since)
    for (let k = 1; k < inside.length; k++) {
      const before = inside[k - 1], after = inside[k]
      const fresh = before.reset && (before.reset <= after.at || after.reset - before.reset > hour)
      drained += (fresh ? after.used : Math.max(0, after.used - before.used)) * capacityOf(accounts[index])
    }
  }
  // Until most of a day is recorded, the burn rate is each account's average since its week began, which needs no
  // history. A week opened in the last few hours counts as six hours old, so its first requests don't read as a steep rate.
  const burn = start - since >= 20 * hour ? drained / ((start - since) / hour) : accounts.reduce((sum, account) => {
    const window = find(account, 'seven_day')
    const reset = resetOf(window)
    return reset ? sum + window.used * capacityOf(account) / Math.max((start - (reset - week)) / hour, 6) : sum
  }, 0)

  const state = accounts.map(account => ({
    capacity: capacityOf(account),
    left: left(find(account, 'seven_day')),
    reset: resetOf(find(account, 'seven_day'))
  }))
  const value = () => state.reduce((sum, account) => sum + account.left * account.capacity, 0)
  const projected = [{ time: start, value: value() }]
  const jumps = []
  let empty = null
  for (let time = start + step; time <= to; time += step) {
    for (const account of state) {
      if (!account.reset || account.reset > time) continue
      jumps.push({ time: account.reset, gain: (100 - account.left) * account.capacity })
      account.left = 100
      account.reset = null
    }
    let need = (burn || 0) * step / hour
    for (const account of state.filter(account => account.left > 0).sort((a, b) => (a.reset ?? Infinity) - (b.reset ?? Infinity))) {
      if (need <= 0) break
      const take = Math.min(account.left * account.capacity, need)
      account.left -= take / account.capacity
      need -= take
      // An untouched week starts with its first use.
      account.reset ??= time + week
      if (account.left < 1e-6) account.left = 0
    }
    if (value() < capacity / 2) empty ??= time
    projected.push({ time, value: value() })
  }
  // Resets within ninety minutes of each other read as one step.
  const resets = []
  for (const jump of jumps.sort((a, b) => a.time - b.time)) {
    const last = resets.at(-1)
    if (last && jump.time - last.time < 1.5 * hour) last.gain += jump.gain
    else resets.push({ ...jump })
  }
  const share = series => series.map(point => ({ time: point.time, value: point.value == null ? null : point.value / capacity }))
  return {
    recorded: share(recorded),
    projected: share(projected),
    resets: resets.map(reset => ({ time: reset.time, gain: reset.gain / capacity })).filter(reset => reset.gain >= 0.5),
    burn: burn / capacity,
    empty
  }
}

// A provider's weekly pool over the last three days and the next eleven: how much is left, when it runs out and when
// it refills.
function forecastCard(name, list) {
  const start = now()
  const from = start - 3 * 86400000, to = start + 11 * 86400000
  const model = forecast(list, from, to)
  if (!model) return null
  const color = providerColors[name]
  const x = time => (time - from) / (to - from) * 100
  const y = value => 100 - value
  // The headline takes the place of the bars' pool numbers, beside the provider's name.
  const total = element('div', 'pool-number', element('b', '', `${Math.round(model.projected[0].value)}`), element('span', 'percent', '%'), element('span', 'left-label', 'left'))
  const summary = element('div', 'forecast-summary flips',
    element('div', 'forecast-total', total, element('span', 'metric-label', 'Weekly pool')))
  flipTarget(total, 'Show the limit bars', name)
  summary.addEventListener('click', event => { if (!event.target.closest('.pool-number')) flipView(name) })

  const plot = element('div', 'plot forecast-plot')
  for (const share of [0, .5, 1]) {
    const line = element('div', 'gridline', element('span', '', `${share * 100}%`))
    line.style.bottom = `${share * 100}%`
    plot.append(line)
  }
  const svgNS = 'http://www.w3.org/2000/svg'
  const svg = document.createElementNS(svgNS, 'svg')
  svg.setAttribute('viewBox', '0 0 100 100')
  svg.setAttribute('preserveAspectRatio', 'none')
  svg.setAttribute('aria-hidden', 'true')
  const defs = document.createElementNS(svgNS, 'defs')
  svg.append(defs)
  const path = (className, d, paint) => {
    const node = document.createElementNS(svgNS, 'path')
    node.setAttribute('class', className)
    node.setAttribute('d', d)
    for (const [key, value] of Object.entries(paint)) node.setAttribute(key, value)
    return node
  }
  for (const [id, opacity] of [['past', .28], ['ahead', .12]]) {
    const fade = document.createElementNS(svgNS, 'linearGradient')
    fade.id = `${id}-${name}`
    fade.setAttribute('x2', '0')
    fade.setAttribute('y2', '1')
    for (const [offset, value] of [[0, opacity], [1, 0]]) {
      const stop = document.createElementNS(svgNS, 'stop')
      stop.setAttribute('offset', offset)
      stop.setAttribute('stop-color', color)
      stop.setAttribute('stop-opacity', value)
      fade.append(stop)
    }
    defs.append(fade)
  }
  const at = point => `${+x(point.time).toFixed(2)},${+y(point.value).toFixed(2)}`
  const shape = (points, kind) => {
    const line = `M${points.map(at).join(' L')}`
    svg.append(path('farea', `${line} L${x(points.at(-1).time)},100 L${x(points[0].time)},100 Z`, { fill: `url(#${kind}-${name})` }),
      path(kind === 'past' ? 'fline' : 'fline projected', line, { stroke: color }))
  }
  // The recorded line breaks at gaps in the history.
  model.recorded.reduce((runs, point) => {
    if (point.value == null) runs.push([])
    else runs.at(-1).push(point)
    return runs
  }, [[]]).filter(run => run.length > 1).forEach(run => shape(run, 'past'))
  // Without a burn rate there is nothing to project.
  shape(model.projected, 'ahead')
  plot.append(svg)

  const marker = (className, time, value, text) => {
    const node = element('div', className, text && element('span', '', text))
    node.style.left = `${x(time)}%`
    if (value != null) node.style.top = `${y(value)}%`
    plot.append(node)
    return node
  }
  marker('now-line', start, null, 'Now')
  // Each reset is labelled above the top of its step, once there is a projection to step. Where labels would
  // overlap, the bigger reset keeps its label.
  const placed = []
  for (const reset of [...model.resets].sort((a, b) => b.gain - a.gain)) {
    const step = model.projected.filter(point => point.time >= reset.time && point.time < reset.time + 1.5 * hour)
    if (!step.length || placed.some(time => Math.abs(x(time) - x(reset.time)) < 5)) continue
    placed.push(reset.time)
    marker('reset-label', reset.time, Math.max(...step.map(point => point.value)), `+${Math.round(reset.gain)}%`)
  }
  if (model.empty) marker('empty-mark', model.empty, 0, `Empty ${when(model.empty)}`)
  marker('now-dot', start, model.projected[0].value).style.setProperty('--c', color)

  // The hovered time reads out above the plot, with a dot on the line.
  const guide = marker('guide', start, null, ' ')
  const reading = guide.firstChild
  reading.setAttribute('aria-live', 'polite')
  const dot = marker('hover-dot', start, 0)
  dot.style.setProperty('--c', color)
  const valueAt = time => {
    const series = time <= start ? model.recorded : model.projected
    return series.findLast(point => point.time <= time + 1)?.value ?? null
  }
  let cursor = null
  const readout = time => {
    cursor = time
    plot.classList.toggle('hover', time != null)
    if (time == null) return
    const value = valueAt(time)
    guide.style.left = dot.style.left = `${x(time)}%`
    dot.hidden = value == null
    if (value != null) dot.style.top = `${y(value)}%`
    // Centred mid-plot, sliding to left- or right-aligned towards the edges.
    reading.style.translate = `${-(Math.min(Math.max(x(time), 15), 85) - 15) / 0.7}% 0`
    const date = new Date(time).toLocaleString(undefined, { weekday: 'short', hour: 'numeric', minute: '2-digit' })
    reading.replaceChildren(element('span', 'reading-time', date), value == null ? 'No reading' : `${Math.round(value)}% left${time > start + hour / 4 ? ', projected' : ''}`)
  }
  plot.addEventListener('pointermove', event => {
    const box = plot.getBoundingClientRect()
    readout(from + Math.min(Math.max((event.clientX - box.left) / box.width, 0), 1) * (to - from))
  })
  plot.addEventListener('pointerleave', () => readout(null))
  plot.tabIndex = 0
  plot.setAttribute('role', 'group')
  plot.setAttribute('aria-label', `${names[name]} weekly pool over time. Use the arrow keys to read it.`)
  plot.addEventListener('keydown', event => {
    const moves = { ArrowLeft: -6 * hour, ArrowRight: 6 * hour }
    if (event.key in moves) readout(Math.min(Math.max((cursor ?? start) + moves[event.key], from), to))
    else if (event.key === 'Home') readout(from)
    else if (event.key === 'End') readout(to)
    else if (event.key === 'Escape') readout(null)
    else return
    event.preventDefault()
  })
  plot.addEventListener('blur', () => readout(null))

  // A label under each day, and a faint line at each midnight. Narrower cards label every other day, then every fourth.
  const axis = element('div', 'xaxis days')
  const midnight = new Date(from)
  midnight.setHours(0, 0, 0, 0)
  for (let index = 0; midnight.getTime() < to; index++, midnight.setDate(midnight.getDate() + 1)) {
    const noon = midnight.getTime() + 12 * hour
    if (midnight.getTime() > from) marker('day-line', midnight.getTime())
    if (x(noon) < 2 || x(noon) > 98) continue
    const label = element('span', index % 4 ? index % 2 ? 'odd' : 'even' : '', new Date(noon).toLocaleDateString(undefined, { weekday: 'short', day: 'numeric' }))
    label.style.left = `${x(noon)}%`
    axis.append(label)
  }
  return { summary, body: element('div', 'forecast-body', plot, axis) }
}

function provider(name, list, owner, models) {
  const usable = poolAccounts(list)
  const logo = element('span', 'provider-mark')
  logo.innerHTML = icons[name]
  const heading = element('h2', '', names[name])
  const add = iconButton(`Add ${names[name]} account`, plusIcon, 'join-button', async () => {
    // A sign-in session lasts five minutes, so reopening soon after keeps the same one.
    if (view.joining?.provider !== name || Date.now() - view.joining.started > 270000) {
      const result = await call('/proxy/api/join/start', { provider: name }, '')
      view.joining = { provider: name, session: result.session, url: result.url, started: Date.now() }
    }
    anchorJoin()
    $('#join').replaceChildren(joinPanel())
    if (!$('#join').matches(':popover-open')) $('#join').showPopover()
  })
  add.dataset.provider = name
  heading.append(add)
  if (owner) {
    const gear = element('button', 'icon-button settings-gear')
    gear.type = 'button'
    gear.dataset.provider = name
    setIcon(gear, `${names[name]} settings`, gearIcon)
    gear.popoverTargetElement = $('#provider-settings')
    // Another provider's gear switches the open panel instead of closing it.
    gear.addEventListener('click', () => {
      gear.popoverTargetAction = view.settingsFor === name ? 'toggle' : 'show'
      view.settingsFor = name
      fillSettings()
    })
    heading.append(gear)
  }
  // An account whose week is spent has no session to use, so its 5-hour window counts as empty.
  const sessions = usable.filter(account => find(account, 'five_hour')).map(account => {
    const week = find(account, 'seven_day')
    return week && left(week) === 0 ? { ...account, windows: account.windows.map(window => window.name === 'five_hour' ? { ...window, used: 100 } : window) } : account
  })
  const sessionPool = poolWindow(sessions, 'five_hour')
  const head = element('div', 'provider-summary',
    element('div', 'provider-identity', logo, heading),
    poolCell('weekly', poolWindow(usable, 'seven_day'), poolShares(usable)),
    ...models.map(model => poolCell(`${familyOf(model)} weekly`, poolWindow(usable.filter(account => find(account, model)), model), [])),
    poolCell('session', sessionPool, []))
  // The forecast replaces the bars while it is shown. Clicking the pool numbers opens it and its headline closes it.
  const card = view.forecasting ? forecastCard(name, list) : null
  if (card) head.replaceChildren(head.firstChild, card.summary)
  else {
    for (const number of head.querySelectorAll('.pool-number')) flipTarget(number, 'Show the weekly forecast', name)
    head.addEventListener('click', event => {
      if (event.target.closest('.pool') && !event.target.closest('button, [data-member], .pool-number')) flipView(name)
    })
    head.classList.add('flips')
  }
  const section = element('section', models.length ? 'provider model-limits' : 'provider', head, ...(card ? [card.body] : list.map(account => accountRow(account, owner, models))))
  section.style.viewTransitionName = `provider-${name}`
  return section
}

// Makes a headline number switch between the bars and the forecast, by pointer or keyboard.
function flipTarget(node, label, name) {
  node.tabIndex = 0
  node.role = 'button'
  node.setAttribute('aria-label', label)
  node.dataset.tip = label
  node.addEventListener('click', () => flipView(name))
  node.addEventListener('keydown', event => {
    if (event.key !== 'Enter' && event.key !== ' ') return
    event.preventDefault()
    flipView(name)
  })
}

// Switches every provider between the limit bars and the weekly forecast, crossfading the sections while the new bars
// sweep in and the forecast lines draw. Focus returns to the same provider's headline.
function flipView(name) {
  view.forecasting = !view.forecasting
  localStorage.setItem('limitsView', view.forecasting ? 'forecast' : 'bars')
  const swap = () => {
    render()
    enter($('#providers'))
    document.querySelector(`#providers .provider:nth-child(${['claude', 'codex'].indexOf(name) + 1}) :is(.pool-number, .forecast-total)`)?.focus({ preventScroll: true })
  }
  if (document.startViewTransition && !matchMedia('(prefers-reduced-motion: reduce)').matches) document.startViewTransition(swap)
  else swap()
}

function niceTop(max) {
  const raw = max / 3
  const power = 10 ** Math.floor(Math.log10(raw))
  const step = [1, 2, 2.5, 5, 10].map(factor => factor * power).find(step => step >= raw)
  return { step, top: Math.ceil(max / step) * step }
}

// Monotone cubic interpolation, so the curve never dips below zero or overshoots a day's value.
function curve(points) {
  if (points.length === 1) points = [{ x: 0, y: points[0].y }, { x: 100, y: points[0].y }]
  const slopes = points.slice(1).map((point, index) => (point.y - points[index].y) / (point.x - points[index].x))
  const tangents = points.map((_, index) => index === 0 ? slopes[0] : index === slopes.length ? slopes[index - 1]
    : slopes[index - 1] * slopes[index] <= 0 ? 0 : (slopes[index - 1] + slopes[index]) / 2)
  slopes.forEach((slope, index) => {
    if (!slope) { tangents[index] = tangents[index + 1] = 0; return }
    const a = tangents[index] / slope, b = tangents[index + 1] / slope, size = Math.hypot(a, b)
    if (size > 3) { tangents[index] = 3 / size * a * slope; tangents[index + 1] = 3 / size * b * slope }
  })
  const at = (x, y) => `${+x.toFixed(2)},${+y.toFixed(2)}`
  return points.slice(1).reduce((path, point, index) => {
    const from = points[index], dx = (point.x - from.x) / 3
    return `${path} C${at(from.x + dx, from.y + tangents[index] * dx)} ${at(point.x - dx, point.y - tangents[index + 1] * dx)} ${at(point.x, point.y)}`
  }, `M${at(points[0].x, points[0].y)}`)
}

const providerColors = { claude: '#d97757', codex: '#f5f5f7' }
const allIcon = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m12 4-9 4.5 9 4.5 9-4.5L12 4ZM3 13l9 4.5 9-4.5M3 17.5 12 22l9-4.5"/></svg>'

// Each member's daily tokens since the pool's first use, on both providers or one.
function chart(members) {
  const filter = view.chartProvider
  const length = members[0]?.days?.claude?.length || 1
  const providers = filter === 'all' ? Object.keys(providerColors) : [filter]
  const series = members.map(member => {
    const values = Array.from({ length }, (_, index) => providers.reduce((sum, provider) => sum + (member.days?.[provider]?.[index] || 0), 0))
    return { name: member.name, values, total: values.reduce((a, b) => a + b, 0) }
  }).filter(item => item.total).sort((a, b) => b.total - a.total)
  const filters = element('div', 'filters', ...['all', ...Object.keys(providerColors)].map(provider => {
    const choice = element('button')
    const label = provider === 'all' ? 'All providers' : names[provider]
    choice.innerHTML = provider === 'all' ? allIcon : icons[provider]
    choice.setAttribute('aria-label', label)
    choice.type = 'button'
    choice.dataset.provider = provider
    choice.setAttribute('aria-pressed', provider === filter)
    choice.addEventListener('click', () => {
      view.chartProvider = provider
      // Members slide to their new order while the chart redraws from the left.
      const swap = () => {
        drawChart()
        drawMembers()
        enter($('#chart'))
        $(`#chart [data-provider="${provider}"]`).focus()
      }
      if (document.startViewTransition && !matchMedia('(prefers-reduced-motion: reduce)').matches) document.startViewTransition(swap)
      else swap()
    })
    return choice
  }))
  const label = element('span', 'chart-label', 'Daily tokens')
  const head = element('div', 'chart-head', label, filters)
  if (!series.length) return [head, element('p', 'muted small', filter === 'all' ? 'No usage yet' : `No ${names[filter]} usage yet`)]
  const { step, top } = niceTop(Math.max(...series.flatMap(item => item.values)))
  const plot = element('div', 'plot')
  for (let value = 0; value <= top + step / 2; value += step) {
    const line = element('div', 'gridline', element('span', '', value ? compact.format(value) : '0'))
    line.style.bottom = `${value / top * 100}%`
    plot.append(line)
  }
  const span = Math.max(length - 1, 1)
  const x = index => length === 1 ? 50 : index / span * 100
  const svgNS = 'http://www.w3.org/2000/svg'
  const svg = document.createElementNS(svgNS, 'svg')
  const defs = document.createElementNS(svgNS, 'defs')
  svg.append(defs)
  svg.setAttribute('viewBox', '0 0 100 100')
  svg.setAttribute('preserveAspectRatio', 'none')
  svg.setAttribute('aria-hidden', 'true')
  const path = (className, name, d, paint) => {
    const node = document.createElementNS(svgNS, 'path')
    node.setAttribute('class', className)
    node.setAttribute('d', d)
    node.dataset.member = name
    for (const [key, value] of Object.entries(paint)) node.setAttribute(key, value)
    return node
  }
  // The heaviest member is painted first, so lighter lines stay visible on top.
  series.forEach(({ name, values }, index) => {
    const line = curve(values.map((value, index) => ({ x: x(index), y: 100 - value / top * 100 })))
    // Each area fades from its line to the baseline.
    const fade = document.createElementNS(svgNS, 'linearGradient')
    fade.id = `fade-${index}`
    fade.setAttribute('x2', '0')
    fade.setAttribute('y2', '1')
    for (const [offset, opacity] of [[0, .22], [1, 0]]) {
      const stop = document.createElementNS(svgNS, 'stop')
      stop.setAttribute('offset', offset)
      stop.setAttribute('stop-color', colorOf(name))
      stop.setAttribute('stop-opacity', opacity)
      fade.append(stop)
    }
    defs.append(fade)
    svg.append(path('area', name, `${line} L100,100 L0,100 Z`, { fill: `url(#fade-${index})` }),
      path('line', name, line, { fill: 'none', stroke: colorOf(name), 'stroke-width': '2', 'vector-effect': 'non-scaling-stroke' }))
  })
  plot.append(svg)
  const last = Date.parse(`${view.data.today}T00:00:00Z`)
  const dayName = index => index === length - 1 ? 'Today' : new Date(last - (length - 1 - index) * 86400000).toLocaleDateString(undefined, { month: 'short', day: 'numeric', timeZone: 'UTC' })
  // The hovered day reads out in the header, and each member row shows its tokens for that day.
  const byName = Object.fromEntries(series.map(item => [item.name, item.values]))
  // A day's tokens beside a mark: the layers icon for both providers together, or a provider's own.
  const tokens = (icon, label, value) => {
    const mark = element('span', 'day-split')
    mark.innerHTML = icon
    mark.append(compact.format(value))
    mark.setAttribute('aria-label', `${label} ${compact.format(value)}`)
    return mark
  }
  const byMember = Object.fromEntries(members.map(member => [member.name, member]))
  const providerDay = (list, provider, index) => list.reduce((sum, member) => sum + (member?.days?.[provider]?.[index] || 0), 0)
  // A member's provider bar shows the hovered day's split, and their whole history again afterwards.
  const daySplit = (row, index) => {
    const values = Object.fromEntries(providers.map(provider => [provider, index == null ? null : providerDay([byMember[row.dataset.person]], provider, index)]))
    const total = Object.values(values).reduce((sum, value) => sum + value, 0)
    for (const part of row.querySelectorAll('.mix i')) {
      const share = index == null ? +part.dataset.share : total ? values[part.dataset.provider] / total * 100 : 0
      part.style.width = `${share}%`
      part.hidden = !share
    }
  }
  const showDay = (index, total) => {
    label.classList.toggle('on', index != null)
    if (index == null) label.textContent = 'Daily tokens'
    else if (filter !== 'all') label.replaceChildren(element('span', '', dayName(index)), tokens(icons[filter], names[filter], total))
    else label.replaceChildren(element('span', '', dayName(index)), tokens(allIcon, 'All providers', total),
      ...providers.map(provider => tokens(icons[provider], names[provider], providerDay(members, provider, index))))
    for (const row of document.querySelectorAll('#members .member')) {
      const day = row.querySelector('.member-total > span')
      day.textContent = index == null ? day.dataset.today : `${compact.format(byName[row.dataset.person]?.[index] || 0)} ${index === length - 1 ? 'today' : `on ${dayName(index)}`}`
      daySplit(row, index)
    }
  }
  // Each hover column is centred on its day's point, so the edge columns overhang the plot by half.
  const cols = element('div', 'cols')
  if (length > 1) cols.style.inset = `0 ${-50 / span}%`
  for (let index = 0; index < length; index++) {
    const col = element('div', 'col')
    const used = series.filter(item => item.values[index]).sort((a, b) => b.values[index] - a.values[index])
    const total = used.reduce((sum, item) => sum + item.values[index], 0)
    col.tabIndex = 0
    for (const type of ['pointerenter', 'focus']) col.addEventListener(type, () => showDay(index, total))
    for (const type of ['pointerleave', 'blur']) col.addEventListener(type, () => showDay(null))
    col.setAttribute('aria-label', `${dayName(index)}: ${compact.format(total)} tokens`)
    cols.append(col)
  }
  plot.append(cols)
  plot.setAttribute('role', 'group')
  plot.setAttribute('aria-label', `Daily tokens per member${filter === 'all' ? '' : ` on ${names[filter]}`}`)
  const axis = element('div', 'xaxis')
  for (const index of new Set([0, Math.floor(span / 2), length - 1])) {
    const label = element('span', '', dayName(index))
    label.style.left = `${x(index)}%`
    axis.append(label)
  }
  if (length > 1) {
    axis.firstChild.style.left = '0'
    axis.lastChild.style.left = '100%'
  }
  return [head, plot, axis]
}

const drawChart = () => $('#chart').replaceChildren(...chart(view.data.members))
const drawMembers = () => {
  const totalFor = member => view.chartProvider === 'all'
    ? member.total
    : (member.days?.[view.chartProvider] || []).reduce((sum, value) => sum + value, 0)
  const sorted = [...view.data.members].sort((a, b) => totalFor(b) - totalFor(a) || (b.last_used || 0) - (a.last_used || 0))
  // Replaced rows never fire pointerleave, so a highlight from the old rows is cleared here.
  focusMember(null)
  $('#members').replaceChildren(...sorted.map(member => memberRow(member, view.data.role === 'owner', view.chartProvider)))
}

// Dims every part of a bar or chart except the named member's.
function highlight(container, name) {
  if (!container) return
  container.classList.toggle('hl', !!name)
  for (const part of container.querySelectorAll('[data-member]')) part.classList.toggle('on', part.dataset.member === name)
}

function memberRow(member, owner, provider = 'all') {
  const selectedDays = provider === 'all' ? null : member.days?.[provider] || []
  const total = selectedDays ? selectedDays.reduce((sum, value) => sum + value, 0) : member.total
  const today = selectedDays ? selectedDays.at(-1) || 0 : member.today_total
  const lastUsed = selectedDays ? member.last_used_by_provider?.[provider] : member.last_used
  const mix = element('div', 'track mix')
  // Every shown provider gets a part, so hovering a day can resize them to that day's split. Unused ones are hidden.
  for (const name of provider === 'all' ? Object.keys(providerColors) : [provider]) {
    const used = member.days?.[name]?.reduce((a, b) => a + b, 0) || 0
    const part = element('i')
    part.dataset.provider = name
    part.dataset.share = total ? used / total * 100 : 0
    part.style.width = `${part.dataset.share}%`
    part.hidden = !used
    part.style.background = providerColors[name]
    mix.append(part)
  }
  const editable = owner && member.managed
  const name = element(editable ? 'button' : 'span', editable ? 'who editable' : 'who', member.name)
  name.style.color = colorOf(member.name)
  const title = element('div', 'member-title', name, member.paused ? element('span', 'badge', 'Paused') : null)
  const day =element('span', 'muted', `${compact.format(today)} today`)
  day.dataset.today = day.textContent
  const totals = element('div', 'member-total', element('b', 'num', compact.format(total)), day)
  // Like an account's recent members, the providers used within ten minutes of the member's last request.
  const recent = (provider === 'all' ? Object.keys(providerColors) : [provider])
    .filter(name => lastUsed && member.last_used_by_provider?.[name] >= lastUsed - 600)
  const time = element('div', 'member-time muted', ...recent.map(name => {
    const mark = element('span', 'used-on')
    mark.innerHTML = icons[name]
    mark.role = 'img'
    mark.setAttribute('aria-label', names[name])
    return mark
  }), element('span', '', ago(lastUsed)))
  // When and where the member was last active reads as part of their name line.
  title.append(time)
  const node = element('div', 'member', element('div', 'member-left', title, mix), totals)
  node.dataset.person = member.name
  node.style.viewTransitionName = `member-${view.data.members.indexOf(member)}`
  node.tabIndex = 0
  if (editable) {
    name.type = 'button'
    name.setAttribute('aria-label', `Rename member ${member.name}`)
    name.addEventListener('click', () => editName(member.name, member.name, title,
      newName => call('/proxy/api/clients', { action: 'rename', name: member.name, new_name: newName }), true))
    // The owner's actions appear after the last-used time while the row is hovered or focused.
    title.append(element('div', 'row-actions',
      iconButton(`${member.paused ? 'Resume' : 'Pause'} ${member.name}`, member.paused ? playIcon : pauseIcon, '', async () => {
        await call('/proxy/api/clients', { action: member.paused ? 'resume' : 'pause', name: member.name })
        await load()
      }),
      iconButton(`Show key for ${member.name}`, 'M15.75 5.25a3 3 0 0 1 3 3m3 0a6 6 0 0 1-7.029 5.912c-.563-.097-1.159.026-1.563.43L10.5 16.5H8.25v2.25H6V21H3v-3l6.908-6.908c.404-.404.527-1 .43-1.563A6 6 0 1 1 21.75 8.25Z', '',
        async () => showPrompt(member.name, (await call('/proxy/api/clients', { action: 'show', name: member.name })).prompt)),
      iconButton(`Revoke key for ${member.name}`, trashIcon, 'danger', async () => {
        if (!await ask(`Revoke ${member.name}'s key?`, 'Their setup stops working immediately.', 'Revoke', true)) return
        await call('/proxy/api/clients', { action: 'revoke', name: member.name })
        await load()
      })))
  }
  return memberHover(node, member.name)
}

const plusIcon = 'M12 5v14M5 12h14'
const pauseIcon = 'M9 5v14M15 5v14'
const playIcon = 'M7 4.5v15l12-7.5Z'
const calendarIcon ='M8 2v4M16 2v4M3 10h18M5 4h14a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2Z'
const trashIcon ='M4 7h16M10 3h4l2 4M6 7l1 14h10l1-14M10 11v6M14 11v6'
const gearIcon = 'M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2zM15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z'

function iconButton(label, path, className, run) {
  const node = button('', `icon-button ${className}`, run)
  setIcon(node, label, path)
  return node
}

function setIcon(node, label, path) {
  node.setAttribute('aria-label', label)
  node.innerHTML = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="${path}"/></svg>`
}

// Posts an owner control's new state, reloads and keeps keyboard focus on the same control.
async function saveControl(control, path, body, revert) {
  const focused = document.activeElement === control
  control.disabled = true
  try {
    await call(path, body)
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
  return element('label', model.enabled ? 'setting' : 'setting off', element('span', 'setting-name', model.name || model.id),
    ownerSwitch(`Allow ${model.id}`, model.enabled, `model:${model.id}`, '/proxy/api/models', enabled => ({ id: model.id, enabled })))
}

const effortNames = { minimal: 'Minimal', low: 'Low', medium: 'Medium', high: 'High', xhigh: 'XHigh', max: 'Max' }

// Effort levels a provider allows. Requests at disabled efforts use the nearest allowed level below.
function effortChips(provider, options) {
  return options.efforts.map(({ level, enabled }) => {
    const node = element('button', enabled ? 'chip on' : 'chip', effortNames[level] || level)
    node.type = 'button'
    node.dataset.key = `option:${provider}:${level}`
    node.setAttribute('aria-pressed', enabled)
    node.addEventListener('click', () => {
      node.classList.toggle('on')
      saveControl(node, '/proxy/api/options', { provider, setting: level, enabled: !enabled }, () => node.classList.toggle('on'))
    })
    return node
  })
}

// Current models, then older versions behind a toggle, as T3 Code does. The toggle sits above the
// older models, so opening them never moves it.
function modelGroup(provider, models) {
  const panel = element('div', 'model-list', ...models.filter(model => !model.older).map(modelRow))
  const legacy = models.filter(model => model.older)
  if (!legacy.length) return panel
  const on = legacy.filter(model => model.enabled).length
  const toggle = element('button', 'legacy-toggle', element('span', '', 'Older models'),
    element('span', 'muted', on ? `${on} of ${legacy.length} on` : `${legacy.length}`))
  toggle.type = 'button'
  toggle.insertAdjacentHTML('beforeend', '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m9 6 6 6-6 6"/></svg>')
  const fold = element('div', 'legacy-fold', element('div', '', element('div', 'model-list', ...legacy.map(modelRow))))
  // Toggled in place rather than re-rendered, so the older models fold open.
  const show = open => {
    panel.classList.toggle('legacy-open', open)
    toggle.setAttribute('aria-expanded', open)
    fold.inert = !open
  }
  toggle.addEventListener('click', () => {
    const open = !view.legacyOpen.has(provider)
    view.legacyOpen[open ? 'add' : 'delete'](provider)
    show(open)
  })
  show(view.legacyOpen.has(provider))
  panel.append(toggle, fold)
  return panel
}

// Fills the settings popover for the provider whose gear opened it, anchored to that gear.
function fillSettings() {
  const panel = $('#provider-settings')
  const name = view.data?.role === 'owner' ? view.settingsFor : null
  for (const gear of document.querySelectorAll('.settings-gear')) {
    gear.classList.toggle('open', gear.dataset.provider === name)
    // Closing keeps the last anchor, so the popover fades out in place instead of jumping to the corner.
    if (name) gear.style.anchorName = gear.dataset.provider === name ? '--settings' : ''
  }
  if (!name) {
    if (panel.matches(':popover-open')) panel.hidePopover()
    return
  }
  const options = view.data.options?.[name]
  panel.replaceChildren(...(options ? [
    element('section', 'setting-group', element('h3', '', 'Reasoning effort'), element('div', 'segments', ...effortChips(name, options))),
    element('section', 'setting-group', element('label', options.fast ? 'setting' : 'setting off', element('span', 'setting-name', 'Fast mode'),
      ownerSwitch(`Allow fast mode on ${names[name]}`, options.fast, `option:${name}:fast`, '/proxy/api/options', enabled => ({ provider: name, setting: 'fast', enabled }))))] : []),
    element('section', 'setting-group', element('h3', '', 'Models'), modelGroup(name, (view.data.models || []).filter(model => model.provider === name))))
}

// Entrance animations play on first load and on filter changes, not on the minute refresh.
function enter(node) {
  node.classList.add('enter')
  clearTimeout(node.entering)
  node.entering = setTimeout(() => node.classList.remove('enter'), 1200)
}

function render() {
  hideTip()
  if (!view.shown) enter($('.board'))
  view.shown = true
  const { accounts, role } = view.data
  const owner = role === 'owner'
  const checked = Math.max(...accounts.map(account => account.checked_at ? Date.parse(account.checked_at) : 0))
  $('#updated').textContent = checked > 0 ? `Limits checked ${ago(checked / 1000)}` : ''
  const byProvider = provider => accounts.filter(account => account.provider === provider)
  $('#providers').replaceChildren(...['claude', 'codex'].map(name => provider(name, byProvider(name), owner, modelWindows(byProvider(name)))))
  drawChart()
  drawMembers()
  $('#add-member').hidden = !owner
  fillSettings()
  anchorJoin()
  if (!owner && $('#add').matches(':popover-open')) $('#add').hidePopover()
  // Signed out, the person button opens the sign-in popover. Signed in, it signs out.
  $('#admin').popoverTargetElement = owner ? null : $('#login')
  setIcon($('#admin'), owner ? 'Sign out' : 'Sign in', owner
    ? 'M9 5H5v14h4m5-14 7 7-7 7M21 12H9'
    : 'M20 21v-2a7 7 0 0 0-14 0v2M12 3a4 4 0 1 0 0 8 4 4 0 0 0 0-8Z')
}

// The sign-in steps for a new account, shown in a popover under the provider's add button. Sign-in runs on the
// server, so it works without a key or port forwarding.
function joinPanel() {
  const joining = view.joining
  const link = element('a', 'join-link', `Open ${names[joining.provider]} sign-in ↗`)
  link.href = joining.url
  link.target = '_blank'
  link.rel = 'noopener noreferrer'
  const field = element('input')
  field.spellcheck = false
  field.autocomplete = 'off'
  field.placeholder = 'Final page address or code'
  field.setAttribute('aria-label', 'Callback URL or authorization code')
  const note = element('p', 'join-status', '')
  note.setAttribute('role', 'status')
  const form = element('form', 'join-form', element('h3', '', `Add a ${names[joining.provider]} account`), link,
    element('p', 'muted', `After signing in, paste the final page's address here, even if it fails to load, or the code ${names[joining.provider]} shows. Press Enter within five minutes.`),
    field, note)
  form.addEventListener('submit', async event => {
    event.preventDefault()
    field.disabled = true
    note.classList.remove('bad')
    try {
      await call('/proxy/api/join/complete', { callback: field.value }, joining.session)
      note.textContent = 'Completing sign-in…'
      for (const deadline = performance.now() + 60000; ;) {
        if (performance.now() > deadline) throw new Error('Sign-in is taking longer than expected. Start again if it does not complete.')
        await new Promise(resolve => setTimeout(resolve, 2000))
        const result = await call('/proxy/api/join/status', {}, joining.session)
        if (result.status === 'error') throw new Error(result.error)
        if (result.status === 'ok') break
      }
      view.joining = null
      note.textContent = `Connected. The account is now part of the pool. Its email shows on this page, and if paid extra usage is on, shared traffic may incur charges.`
      await load()
    } catch (error) {
      note.textContent = error.message
      note.classList.add('bad')
      field.disabled = false
      field.focus()
    }
  })
  return form
}

// Keeps the join popover anchored to its provider's add button, which each render rebuilds.
function anchorJoin() {
  if (!view.joining) return
  for (const add of document.querySelectorAll('.join-button')) add.style.anchorName = add.dataset.provider === view.joining.provider ? '--join' : ''
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

const tip = $('#tip')
let tipTarget = null
let hiding = 0
let opening = 0
let closed = 0
function hideTip() {
  clearTimeout(hiding)
  clearTimeout(opening)
  if (tip.matches(':popover-open')) { tip.hidePopover(); closed = performance.now() }
  tipTarget = null
}
function showTip(target) {
  clearTimeout(hiding)
  clearTimeout(opening)
  if (target === tipTarget) return
  const content = target.tip ? target.tip() : target.dataset.tip
  if (!content || !content.length) return hideTip()
  tipTarget = target
  tip.replaceChildren(...[content].flat())
  tip.hidden = false
  if (!tip.matches(':popover-open')) tip.showPopover()
  tip.classList.remove('below')
  const box = target.getBoundingClientRect()
  const width = tip.offsetWidth
  const height = tip.offsetHeight
  const centre = box.left + box.width / 2
  const x = Math.min(Math.max(centre - width / 2, 8), innerWidth - width - 8)
  const below = (target.matches('.icon-button') && box.bottom + height + 12 < innerHeight) || box.top - height - 12 < 8
  const y = below ? box.bottom + 12 : box.top - height - 12
  tip.classList.toggle('below', below)
  tip.style.left = `${x}px`
  tip.style.top = `${y}px`
  tip.style.setProperty('--ax', `${Math.min(Math.max(centre - x, 10), width - 10)}px`)
}
document.addEventListener('pointerover', event => {
  const target = event.target.closest?.('[data-tip]')
  // Tips wait for the pointer to rest, so sweeping across the page shows nothing. Once one is open, or just
  // closed, neighbours show at once.
  if (target && (tipTarget || performance.now() - closed < 300)) showTip(target)
  else if (target) { clearTimeout(opening); opening = setTimeout(() => showTip(target), 500) }
  else {
    clearTimeout(opening)
    // A short grace keeps the tip while the pointer crosses the gap between neighbours.
    if (tipTarget) { clearTimeout(hiding); hiding = setTimeout(hideTip, 120) }
  }
})
document.addEventListener('focusin', event => {
  const target = event.target.closest?.('[data-tip]')
  if (target && target.matches(':focus-visible')) showTip(target)
  else hideTip()
})
document.addEventListener('focusout', hideTip)
addEventListener('scroll', hideTip, { passive: true })



$('#admin').addEventListener('click', () => {
  if (view.data?.role === 'owner') {
    sessionStorage.removeItem('adminKey')
    view.key = ''
    $('#created').close()
    load().catch(report)
  }
})
$('#login').addEventListener('toggle', () => { $('#login-error').textContent = '' })
$('#provider-settings').addEventListener('toggle', event => {
  if (event.newState === 'open') return
  view.settingsFor = null
  fillSettings()
})
$('#login').addEventListener('submit', async event => {
  event.preventDefault()
  view.key = $('#key').value.trim()
  try {
    await load()
    sessionStorage.setItem('adminKey', view.key)
    $('#login').hidePopover()
    $('#key').value = ''
  } catch (error) {
    view.key = ''
    $('#login-error').textContent = error.message
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
    $('#add').hidePopover()
    showPrompt(name, data.prompt)
    await load()
  } catch (error) {
    report(error)
  }
})
setIcon($('#add-member'), 'Add member', plusIcon)
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
  const outside = event => {
    const bounds = dialog.getBoundingClientRect()
    return event.target === dialog && (event.clientX < bounds.left || event.clientX > bounds.right || event.clientY < bounds.top || event.clientY > bounds.bottom)
  }
  // Only a click that starts and ends outside closes it, so a text selection dragged past the edge does not.
  let pressedOutside = false
  dialog.addEventListener('close', hideTip)
  dialog.addEventListener('pointerdown', event => { pressedOutside = outside(event) })
  dialog.addEventListener('click', event => {
    if (pressedOutside && outside(event)) dialog.close()
  })
}
$('#ask-ok').addEventListener('click', () => $('#ask').close('ok'))
$('#ask-cancel').addEventListener('click', () => $('#ask').close())
setIcon($('#refresh'), 'Refresh limits', 'M20 7v5h-5M4 17v-5h5M6 7a7 7 0 0 1 12-1l2 2M4 16l2 2a7 7 0 0 0 12-1')

// Reloading re-sorts each provider's accounts, so priority changes move rows up or down.
// Refreshing waits while a name is being edited or a popover is open, so typing and focus stay put.
// A refresh that succeeds also clears any earlier error or notice.
setInterval(() => {
  if (!document.hidden && !document.querySelector('.name-editor, .pop:popover-open')) load().then(() => { status.textContent = '' }, report)
}, 60000)

load().catch(error => {
  // A rejected admin key falls back to the public view. Other failures keep it for the next refresh.
  if (!view.key || (error.status !== 401 && error.status !== 403)) return report(error)
  sessionStorage.removeItem('adminKey')
  view.key = ''
  load().catch(report)
})
